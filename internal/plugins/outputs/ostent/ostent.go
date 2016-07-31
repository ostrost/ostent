package ostent

import (
	"fmt"
	"sync"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/outputs"
	// "github.com/influxdata/telegraf/plugins/serializers"
)

type group struct {
	mutex sync.Mutex
	kv    map[string]string
}

type ostent struct {
	// serializer serializers.Serializer
	// Metrics map[string]*Metric
	system group
}

func (o *ostent) SystemCopy() map[string]string {
	o.system.mutex.Lock()
	defer o.system.mutex.Unlock()
	copy := make(map[string]string, len(o.system.kv))
	for k, v := range o.system.kv {
		copy[k] = v // v is a string
	}
	return copy
}

func (o *ostent) writeSystem(m telegraf.Metric) error {
	fields := m.Fields()
	for k, field := range fields {
		if k != "uptime_format" { // && !strings.HasPrefix(k, "load")
			continue
		}
		var tail string
		/* if up, ok := fields["uptime"]; ok {
			if uptime, ok := up.(int64); ok {
				tail = fmt.Sprintf(":%02d", 60+uptime%40)
			}
		} // */

		o.system.mutex.Lock()
		defer o.system.mutex.Unlock()
		o.system.kv[k] = fmt.Sprintf("%v", field) + tail
		return nil
	}
	return nil
}

/*
type Metric struct{ value string }
func (m Metric) String() string { return m.value }

func (o *ostent) writeMetric(m telegraf.Metric) error {
	name := m.Name()
	for k, field := range m.Fields() {
		o.Metrics[name+"."+k] = &Metric{
			value: fmt.Sprintf("%#v", field),
		}
	}
	return nil
} // */

func (o *ostent) Connect() error       { return nil }
func (o *ostent) Close() error         { return nil }
func (o *ostent) SampleConfig() string { return `` }
func (o *ostent) Description() string  { return "in-memory output" }

// func (o *ostent) SetSerializer(s serializers.Serializer) { o.serializer = s }

func (o *ostent) Write(ms []telegraf.Metric) error {
	if len(ms) == 0 {
		return nil
	}
	for _, m := range ms {
		if m.Name() == "system" {
			if err := o.writeSystem(m); err != nil {
				return err
			}
		}
		/* if err := o.writeMetric(m); err != nil {
			return err
		} // */
	}
	return nil
}

var Output = &ostent{
	// Metrics: make(map[string]*Metric),
	system: group{kv: make(map[string]string)},
}

func init() { outputs.Add("ostent", func() telegraf.Output { return Output }) }
