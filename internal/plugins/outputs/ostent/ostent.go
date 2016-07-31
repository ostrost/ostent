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
	systemOstent group
}

func (o *ostent) SystemOstentCopy() map[string]string {
	o.systemOstent.mutex.Lock()
	defer o.systemOstent.mutex.Unlock()
	copy := make(map[string]string, len(o.systemOstent.kv))
	for k, v := range o.systemOstent.kv {
		copy[k] = v // v is a string
	}
	return copy
}

func (o *ostent) writeSystemOstent(m telegraf.Metric) error {
	fields := m.Fields()
	// for k, field := range fields { if k = "uptime_format" || strings.HasPrefix(k, "load") //...
	k := "uptime_format"
	if field, ok := fields["uptime_format"]; ok {
		var tail string
		/* if up, ok := fields["uptime"]; ok {
			if uptime, ok := up.(int64); ok {
				tail = fmt.Sprintf(":%02d", 60+uptime%40)
			}
		} // */

		o.systemOstent.mutex.Lock()
		defer o.systemOstent.mutex.Unlock()
		o.systemOstent.kv[k] = fmt.Sprintf("%v", field) + tail
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
		if m.Name() == "system_ostent" {
			if err := o.writeSystemOstent(m); err != nil {
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
	systemOstent: group{kv: make(map[string]string)},
}

func init() { outputs.Add("ostent", func() telegraf.Output { return Output }) }
