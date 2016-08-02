package agent

import (
	"log"
	"runtime"
	"sync"
	"time"

	"github.com/influxdata/telegraf"

	"github.com/ostrost/ostent/internal/config"
	internal_models "github.com/ostrost/ostent/internal/models"
)

func Run(c *config.Config) error {
	if err := c.LoadConfig(); err != nil {
		return err
	}

	a, err := NewAgent(c)
	if err != nil {
		return err
	}

	if err := a.Connect(); err != nil {
		return err
	}
	/* There will be loop with waiting for reload signal.
	reload := make(chan bool, 1)
	reload <- true
	for <-reload {
		reload <- false // */
	if err := a.Run(); err != nil {
		return err
	}
	return nil
}

// Agent runs the agent and collects data based on the given config
type Agent struct {
	Config *config.Config
}

// NewAgent returns an Agent struct based off the given Config
func NewAgent(config *config.Config) (*Agent, error) {
	a := &Agent{
		Config: config,
	}
	return a, nil
}

// Connect connects to all configured outputs
func (a *Agent) Connect() error {
	for _, o := range a.Config.Outputs {
		err := o.Output.Connect()
		if err != nil {
			log.Printf("Failed to connect to output %s, retrying in 15s, "+
				"error was '%s' \n", o.Name, err)
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

func panicRecover(input *internal_models.RunningInput) {
	if err := recover(); err != nil {
		trace := make([]byte, 2048)
		runtime.Stack(trace, true)
		log.Printf("FATAL: Input [%s] panicked: %s, Stack:\n%s\n",
			input.Name, err, trace)
	}
}

// gatherer runs the inputs that have been configured with their own
// reporting interval.
func (a *Agent) gatherer(
	input *internal_models.RunningInput,
	interval time.Duration,
	metricC chan telegraf.Metric,
) error {
	defer panicRecover(input)

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		acc := NewAccumulator(metricC)

		gatherWithTimeout(input, acc, interval)
		select {
		case <-ticker.C:
			continue
		}
	}
}

// gatherWithTimeout gathers from the given input, with the given timeout.
//   when the given timeout is reached, gatherWithTimeout logs an error message
//   but continues waiting for it to return. This is to avoid leaving behind
//   hung processes, and to prevent re-calling the same hung process over and
//   over.
func gatherWithTimeout(
	input *internal_models.RunningInput,
	acc *accumulator,
	timeout time.Duration,
) {
	ticker := time.NewTicker(timeout)
	defer ticker.Stop()
	done := make(chan error)
	go func() {
		done <- input.Input.Gather(acc)
	}()

	for {
		select {
		case err := <-done:
			if err != nil {
				log.Printf("ERROR in input [%s]: %s", input.Name, err)
			}
			return
		case <-ticker.C:
			log.Printf("ERROR: input [%s] took longer to collect than "+
				"collection interval (%s)",
				input.Name, timeout)
			continue
		}
	}
}

// flush writes a list of metrics to all configured outputs
func (a *Agent) flush() {
	var wg sync.WaitGroup

	wg.Add(len(a.Config.Outputs))
	for _, o := range a.Config.Outputs {
		go func(output *internal_models.RunningOutput) {
			defer wg.Done()
			err := output.Write()
			if err != nil {
				log.Printf("Error writing to output [%s]: %s\n",
					output.Name, err.Error())
			}
		}(o)
	}

	wg.Wait()
}

// flusher monitors the metrics input channel and flushes on the minimum interval
func (a *Agent) flusher(metricC chan telegraf.Metric) error {
	// Inelegant, but this sleep is to allow the Gather threads to run, so that
	// the flusher will flush after metrics are collected.
	time.Sleep(time.Millisecond * 200)

	ticker := time.NewTicker(a.Config.Agent.FlushInterval.Duration)

	for {
		select {
		case <-ticker.C:
			a.flush()
		case m := <-metricC:
			for i, o := range a.Config.Outputs {
				if i == len(a.Config.Outputs)-1 {
					o.AddMetric(m)
				} else {
					o.AddMetric(copyMetric(m))
				}
			}
		}
	}
}

func copyMetric(m telegraf.Metric) telegraf.Metric {
	t := time.Time(m.Time())

	tags := make(map[string]string)
	fields := make(map[string]interface{})
	for k, v := range m.Tags() {
		tags[k] = v
	}
	for k, v := range m.Fields() {
		fields[k] = v
	}

	out, _ := telegraf.NewMetric(m.Name(), tags, fields, t)
	return out
}

// Run runs the agent daemon, gathering every Interval
func (a *Agent) Run() error {
	nextDelta := func(d time.Duration) time.Duration {
		now := time.Now()
		return now.Truncate(d).Add(d).Sub(now)
	}

	// channel shared between all input threads for accumulating metrics
	metricC := make(chan telegraf.Metric, 10000)

	time.Sleep(nextDelta(a.Config.Agent.Interval.Duration))

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := a.flusher(metricC); err != nil {
			log.Printf("Flusher routine failed, exiting: %s\n", err.Error())
		}
	}()

	wg.Add(len(a.Config.Inputs))
	for _, input := range a.Config.Inputs {
		interval := a.Config.Agent.Interval.Duration
		// overwrite global interval if this plugin has it's own.
		if input.Config.Interval != 0 {
			interval = input.Config.Interval
		}
		go func(in *internal_models.RunningInput, interv time.Duration) {
			defer wg.Done()
			if err := a.gatherer(in, interv, metricC); err != nil {
				log.Printf(err.Error())
			}
		}(input, interval)
	}

	wg.Wait()
	return nil
}
