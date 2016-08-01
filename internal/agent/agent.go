package agent

import (
	"log"
	"math"
	"sync"
	"time"

	"github.com/influxdata/telegraf"
	_ "github.com/influxdata/telegraf/plugins/outputs/file"

	"github.com/ostrost/ostent/internal/config"
	internal_models "github.com/ostrost/ostent/internal/models"
	_ "github.com/ostrost/ostent/internal/plugins/outputs/ostent" // "ostent" output
	_ "github.com/ostrost/ostent/system_ostent"                   // "system_ostent" input
)

func Start() {
	if err := start(); err != nil {
		log.Printf("Agent error: %s", err)
	}
}

func start() error {
	c := config.NewConfig()
	if err := c.LoadConfig(); err != nil {
		return err
	}

	ag, err := NewAgent(c)
	if err != nil {
		return err
	}

	if err := ag.Connect(); err != nil {
		return err
	}
	/* There will be loop with waiting for reload signal.
	reload := make(chan bool, 1)
	reload <- true
	for <-reload {
		reload <- false // */
	if err := ag.Run(); err != nil {
		return err
	}
	return nil
}

// Agent runs the agent and collects data based on the given config
type Agent struct {
	metricch chan telegraf.Metric
	Config   *config.Config
}

// NewAgent returns an Agent struct based off the given Config
func NewAgent(config *config.Config) (*Agent, error) {
	a := &Agent{
		Config: config,
	}
	a.metricch = make(chan telegraf.Metric, 10000)
	return a, nil
}

func (a *Agent) Connect() error {
	for _, o := range a.Config.Outputs {
		err := o.Output.Connect()
		if err != nil {
			log.Printf("Failed to connect to output %q (will retry in 15s): %s",
				o.Name, err.Error())
			time.Sleep(15 * time.Second)
			err = o.Output.Connect()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Close closes the connection to all configured outputs
func (a *Agent) Close() error {
	var err error
	for _, o := range a.Config.Outputs {
		err = o.Output.Close()
		switch ot := o.Output.(type) {
		case telegraf.ServiceOutput:
			ot.Stop()
		}
	}
	return err
}

func NextDelta(d time.Duration) time.Duration {
	now := time.Now()
	return now.Truncate(d).Add(d).Sub(now)
}

func (ag *Agent) Run() error {
	time.Sleep(NextDelta(ag.Config.Agent.Interval.Duration))

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := ag.flush(); err != nil {
			log.Printf("Flush failed: %s\n", err.Error())
		}
	}()

	wg.Add(len(ag.Config.Inputs))
	for _, input := range ag.Config.Inputs {
		interval := ag.Config.Agent.Interval.Duration
		if input.Config.Interval != 0 {
			interval = input.Config.Interval
		}
		go func(in *internal_models.RunningInput, interv time.Duration, mch chan telegraf.Metric) {
			defer wg.Done()
			if err := gather(in, interv, mch); err != nil {
				log.Printf(err.Error())
			}
		}(input, interval, ag.metricch)
	}
	wg.Wait()
	return nil
}

func (ag *Agent) flush() error {
	time.Sleep(time.Millisecond * 200) // time for gather threads to run

	ticker := time.NewTicker(ag.Config.Agent.FlushInterval.Duration)

	for {
		select {
		case <-ticker.C:
			ag.doFlush()
		case m := <-ag.metricch:
			ag.AddMetric(m)
		}
	}
}

func (ag *Agent) doFlush() {
	var wg sync.WaitGroup
	wg.Add(len(ag.Config.Outputs))
	for _, out := range ag.Config.Outputs {
		go func(output *internal_models.RunningOutput) {
			defer wg.Done()
			if err := output.Write(); err != nil {
				log.Printf("Error writing to output [%s]: %s\n", output.Name, err.Error())
			}
		}(out)
	}
	wg.Wait()
}

func (ag *Agent) AddMetric(m telegraf.Metric) {
	for i, out := range ag.Config.Outputs {
		if i == len(ag.Config.Outputs)-1 { // the last
			out.AddMetric(m)
		} else if mcopy, err := copyMetric(m); err == nil { // err is ignored
			out.AddMetric(mcopy)
		}
	}
}

func copyMetric(m telegraf.Metric) (telegraf.Metric, error) {
	var (
		tags   = make(map[string]string)
		fields = make(map[string]interface{})
	)
	for k, v := range m.Tags() {
		tags[k] = v
	}
	for k, v := range m.Fields() {
		fields[k] = v
	}
	return telegraf.NewMetric(m.Name(), tags, fields, time.Time(m.Time()))
}

type acc struct {
	metricch chan telegraf.Metric
	debug    bool
}

func (ac acc) Debug() bool       { return ac.debug }
func (ac *acc) SetDebug(on bool) { ac.debug = on }

func (ac *acc) SetPrecision(precision, interval time.Duration) {} // TODO
func (ac *acc) DisablePrecision()                              {} // TODO

func (ac *acc) AddError(err error) { log.Printf("Error in input: %s", err) }

// Add of telegraf.Accumulator interface.
func (ac *acc) Add(measurement string, value interface{}, tags map[string]string, t ...time.Time) {
	ac.AddFields(measurement, map[string]interface{}{"value": value}, tags, t...)
}

// AddFields of telegraf.Accumulator interface.
func (ac *acc) AddFields(measurement string, fields map[string]interface{}, tags map[string]string, t ...time.Time) {
	if measurement == "" || len(fields) == 0 {
		return
	}
	if tags == nil {
		tags = make(map[string]string)
	}
	values := make(map[string]interface{})
	for k, v := range fields {

		// Validate uint64 and float64 fields
		switch val := v.(type) {
		case uint64:
			// InfluxDB does not support writing uint64
			if val < uint64(9223372036854775808) {
				values[k] = int64(val)
			} else {
				values[k] = int64(9223372036854775807)
			}
			continue
		case float64:
			// NaNs are invalid values in influxdb, skip measurement
			if math.IsNaN(val) || math.IsInf(val, 0) {
				if false { // if ac.debug TODO
					log.Printf(
						"Measurement [%s] field [%s] has a NaN or Inf field, skipping",
						measurement, k)
				}
				continue
			}
		}

		values[k] = v
	}
	if len(values) == 0 {
		return
	}

	var ts time.Time
	if len(t) > 0 {
		ts = t[0]
	} else {
		ts = time.Now()
	}
	// timestamp = timestamp.Round(ac.precision) // TODO

	m, err := telegraf.NewMetric(measurement, tags, values, ts)
	if err != nil {
		log.Printf("Error adding point [%s]: %s\n", measurement, err.Error())
		return
	}
	// if ac.trace { fmt.Println("> " + m.String()) } // TODO
	ac.metricch <- m
}

func gather(input *internal_models.RunningInput, interval time.Duration, metricch chan telegraf.Metric) error {
	defer catch(input)

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		gatherTimed(input, interval, &acc{metricch: metricch})
		select {
		case <-ticker.C:
			continue
		}
	}
}

func gatherTimed(input *internal_models.RunningInput, interval time.Duration, ac *acc) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	done := make(chan error)
	go func() {
		done <- input.Input.Gather(ac)
	}()

	for {
		select {
		case err := <-done:
			if err != nil {
				log.Printf("Error in input [%s]: %s", input.Name, err)
			}
			return
		case <-ticker.C:
			log.Printf(
				"Error: input [%s] took longer to collect than collection interval (%s)",
				input.Name, interval)
			continue
		}
	}
}

func catch(input *internal_models.RunningInput) {
	if err := recover(); err != nil {
		log.Printf("Panic: %s: %s\n", input.Name, err)
	}
}
