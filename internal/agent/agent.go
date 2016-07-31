package agent

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"strconv"
	"sync"
	"time"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/inputs"
	_ "github.com/influxdata/telegraf/plugins/inputs/system" // "system" input
	"github.com/influxdata/telegraf/plugins/outputs"
	_ "github.com/influxdata/telegraf/plugins/outputs/file"
	"github.com/influxdata/telegraf/plugins/serializers"
	"github.com/influxdata/toml"
	"github.com/influxdata/toml/ast"

	_ "github.com/ostrost/ostent/internal/plugins/outputs/ostent" // "ostent" output
)

func Start() {
	if err := start(); err != nil {
		log.Printf("Agent error: %s", err)
	}
}

func parse(contents []byte) (*ast.Table, error) {
	contents = bytes.TrimPrefix(contents, []byte("\xef\xbb\xbf")) // windows
	for _, dword := range regexp.MustCompile(`\$\w+`).FindAll(contents, -1) {
		if val := os.Getenv(string(dword[1:])); val != "" {
			contents = bytes.Replace(contents, dword, []byte(val), 1)
		}
	}
	return toml.Parse(contents)
}

func (c *config) loadConfig() error {
	tbl, err := parse([]byte(`
[agent]
  interval = "1s"
  flushInterval = "1s"
[[inputs.system]]
  interval = "1s"
# [[outputs.file]]
[[outputs.ostent]]
`))
	if err != nil {
		return err
	}
	if val, ok := tbl.Fields["agent"]; ok {
		subt, ok := val.(*ast.Table)
		if !ok {
			return fmt.Errorf("Cannot parse config")
		}
		if err := toml.UnmarshalTable(subt, c.agent); err != nil {
			return fmt.Errorf("Cannot parse config: [agent] section: %s", err)
		}
	}

	for name, val := range tbl.Fields {
		subt, ok := val.(*ast.Table)
		if !ok {
			return fmt.Errorf("Cannot parse config")
		}
		if name == "outputs" {
			for pname, pval := range subt.Fields {
				switch psubt := pval.(type) {
				case *ast.Table:
					if err := c.addOutput(pname, psubt); err != nil {
						return fmt.Errorf("Parse error: %s", err)
					}
				case []*ast.Table:
					for _, t := range psubt {
						if err := c.addOutput(pname, t); err != nil {
							return fmt.Errorf("Parse error: %s", err)
						}
					}
				default:
					return fmt.Errorf("Unsupported type in config: [%s] section", pname)
				}
			}
		} else if name == "inputs" {
			for pname, pval := range subt.Fields {
				switch psubt := pval.(type) {
				case *ast.Table:
					if err := c.addInput(pname, psubt); err != nil {
						return fmt.Errorf("Parse error: %s", err)
					}
				case []*ast.Table:
					for _, t := range psubt {
						if err := c.addInput(pname, t); err != nil {
							return fmt.Errorf("Parse error: %s", err)
						}
					}
				default:
					return fmt.Errorf("Unsupported type in config: [%s] section", pname)
				}
			}
		}
	}
	return err
}

type duration struct{ Duration time.Duration }

func (d *duration) UnmarshalTOML(b []byte) error {
	var err error
	d.Duration, err = time.ParseDuration(string(b[1 : len(b)-1]))
	if err == nil {
		return nil
	}
	if n, err := strconv.ParseInt(string(b), 10, 64); err == nil {
		d.Duration = time.Second * time.Duration(n)
	} else if f, err := strconv.ParseFloat(string(b), 64); err == nil {
		d.Duration = time.Second * time.Duration(f)
	}
	return nil
}

// fields are public for unmarshalling
type agentConfig struct{ Interval, FlushInterval duration }
type config struct {
	agent *agentConfig

	Inputs  []*input
	Outputs []*output
}

func newConfig() *config {
	return &config{
		agent: &agentConfig{
			// values are defaults
			Interval:      duration{Duration: time.Second * 10},
			FlushInterval: duration{Duration: time.Second * 10},
		},
	}
}

func start() error {
	c := newConfig()
	if err := c.loadConfig(); err != nil {
		return err
	}

	ag := agent{metricch: make(chan telegraf.Metric, 10000), config: c}

	if err := ag.connect(); err != nil {
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

type agent struct {
	metricch chan telegraf.Metric
	config   *config
}

func (ag *agent) connect() error {
	for _, out := range ag.config.Outputs {
		if err := connectOne(out); err != nil {
			return err
		}
	}
	return nil
}

func connectOne(out *output) error {
	err := out.output.Connect()
	if err == nil {
		return nil
	}
	log.Printf("Failed to connect to output %q (will retry in 15s): %s",
		out.name, err.Error())
	time.Sleep(15 * time.Second)
	err = out.output.Connect()
	if err == nil {
		return nil
	}
	log.Printf("Failed to connect to output %q: %s", out.name, err.Error())
	return err
}

func NextDelta(d time.Duration) time.Duration {
	now := time.Now()
	return now.Truncate(d).Add(d).Sub(now)
}

func (ag *agent) Run() error {
	time.Sleep(NextDelta(ag.config.agent.Interval.Duration))

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := ag.flush(); err != nil {
			log.Printf("Flush failed: %s\n", err.Error())
		}
	}()

	wg.Add(len(ag.config.Inputs))
	for _, in := range ag.config.Inputs {
		interval := ag.config.agent.Interval.Duration
		if in.config.interval != 0 {
			interval = in.config.interval
		}
		go func(inputInput *input, inputInterval time.Duration, mch chan telegraf.Metric) {
			defer wg.Done()
			if err := gather(inputInput, inputInterval, mch); err != nil {
				log.Printf(err.Error())
			}
		}(in, interval, ag.metricch)
	}
	wg.Wait()
	return nil
}

func (ag *agent) flush() error {
	time.Sleep(time.Millisecond * 200) // time for gather threads to run

	ticker := time.NewTicker(ag.config.agent.FlushInterval.Duration)

	for {
		select {
		case <-ticker.C:
			ag.doFlush()
		case m := <-ag.metricch:
			ag.addMetric(m)
		}
	}
}

func (ag *agent) doFlush() {
	var wg sync.WaitGroup
	wg.Add(len(ag.config.Outputs))
	for _, out := range ag.config.Outputs {
		go func(o *output) {
			defer wg.Done()
			if err := o.Write(); err != nil {
				log.Printf("Error writing to output [%s]: %s\n", o.name, err.Error())
			}
		}(out)
	}
	wg.Wait()
}

func (ag *agent) addMetric(m telegraf.Metric) {
	for i, out := range ag.config.Outputs {
		if i == len(ag.config.Outputs)-1 { // the last
			out.addMetric(m)
		} else if mcopy, err := copyMetric(m); err == nil { // err is ignored
			out.addMetric(mcopy)
		}
	}
}

func (out *output) addMetric(m telegraf.Metric) {
	out.buf.Add(m)
	if out.buf.Len() == out.bufBatchSize {
		b := out.buf.Batch(out.bufBatchSize)
		if err := out.rawWrite(b); err != nil {
			panic(err)
		}
	}
}

func (out *output) rawWrite(ms []telegraf.Metric) error {
	if ms == nil || len(ms) == 0 {
		return nil
	}
	return out.output.Write(ms)
}

func (out *output) Write() error {
	return out.rawWrite(out.buf.Batch(out.bufBatchSize))
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

type buffer struct{ buf chan telegraf.Metric }

func newBuffer(size int) *buffer { return &buffer{buf: make(chan telegraf.Metric, size)} }

func (b *buffer) Len() int { return len(b.buf) }

func (b *buffer) Add(ms ...telegraf.Metric) {
	for i := range ms {
		select {
		case b.buf <- ms[i]:
		default:
			<-b.buf
			b.buf <- ms[i]
		}
	}
}

func (b *buffer) Batch(bsize int) []telegraf.Metric {
	n := b.Len()
	if n >= bsize {
		n = bsize
	}
	o := make([]telegraf.Metric, n)
	for i := 0; i < n; i++ {
		o[i] = <-b.buf
	}
	return o
}

func makeSerializer(name string) (serializers.Serializer, error) {
	return serializers.NewSerializer(&serializers.Config{
		DataFormat: "graphite",
	})
}

type output struct {
	output telegraf.Output
	name   string
	// Config struct { Interval time.Duration }
	bufBatchSize int
	buf          *buffer
}

func (c *config) addOutput(name string, table *ast.Table) error {
	create, ok := outputs.Outputs[name]
	if !ok {
		return fmt.Errorf("Unknown output by name %q", name)
	}
	out := create()

	if ser, ok := out.(serializers.SerializerOutput); ok {
		newSer, err := makeSerializer(name)
		if err != nil {
			return err
		}
		if newSer == nil {
			return fmt.Errorf("Serializer is nil")
		}
		ser.SetSerializer(newSer)
	}

	bbs := 1000
	c.Outputs = append(c.Outputs, &output{output: out, name: name, bufBatchSize: bbs, buf: newBuffer(bbs)})
	return nil
}

type input struct {
	input  telegraf.Input
	name   string
	config *inputConfig
}

type inputConfig struct {
	name     string
	interval time.Duration
}

func (c *config) addInput(name string, table *ast.Table) error {
	create, ok := inputs.Inputs[name]
	if !ok {
		return fmt.Errorf("Unknown input by name %q", name)
	}
	in := create()
	pc, err := buildInput(name, table)
	if err != nil {
		return err
	}
	c.Inputs = append(c.Inputs, &input{input: in, name: name, config: pc})
	return nil
}

func buildInput(name string, tbl *ast.Table) (*inputConfig, error) {
	conf := &inputConfig{name: name, interval: time.Second * 10}
	if node, ok := tbl.Fields["interval"]; ok {
		if kv, ok := node.(*ast.KeyValue); ok {
			if sv, ok := kv.Value.(*ast.String); ok {
				d, err := time.ParseDuration(sv.Value)
				if err != nil {
					return nil, err
				}
				conf.interval = d
			}
		}
	}
	return conf, nil
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

func gather(input *input, interval time.Duration, metricch chan telegraf.Metric) error {
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

func gatherTimed(input *input, interval time.Duration, ac *acc) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	done := make(chan error)
	go func() {
		done <- input.input.Gather(ac)
	}()

	for {
		select {
		case err := <-done:
			if err != nil {
				log.Printf("Error in input [%s]: %s", input.name, err)
			}
			return
		case <-ticker.C:
			log.Printf(
				"Error: input [%s] took longer to collect than collection interval (%s)",
				input.name, interval)
			continue
		}
	}
}

func catch(input *input) {
	if err := recover(); err != nil {
		log.Printf("Panic: %s: %s\n", input.name, err)
	}
}
