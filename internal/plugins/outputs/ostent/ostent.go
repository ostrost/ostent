package ostent

import (
	"sort"
	"sync"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/outputs"
	// "github.com/influxdata/telegraf/plugins/serializers"
)

type group struct {
	mutex sync.Mutex
	kv    map[string]string
}

type cpu struct {
	N       string
	UserPct int
	SysPct  int
	WaitPct int
	IdlePct int
}

type groupCPU struct {
	mutex sync.Mutex
	list  []cpu
}

type cpuList []cpu

// Len, Swap and Less satisfy sorting interface.
func (cl cpuList) Len() int           { return len(cl) }
func (cl cpuList) Swap(i, j int)      { cl[i], cl[j] = cl[j], cl[i] }
func (cl cpuList) Less(i, j int) bool { return cl[i].IdlePct < cl[j].IdlePct }

type ostent struct {
	// serializer serializers.Serializer
	// Metrics map[string]*Metric
	systemCPU    groupCPU
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

func (o *ostent) SystemCPUCopy() []cpu {
	o.systemCPU.mutex.Lock()
	defer o.systemCPU.mutex.Unlock()
	llen := len(o.systemCPU.list)
	if llen == 0 {
		return []cpu{}
	}
	tshift := 0 // "total shift"
	if o.systemCPU.list[llen-1].N == "cpu-total" {
		tshift = 1
	}

	copy := make([]cpu, llen)
	for i, c := range o.systemCPU.list[:llen-tshift] {
		copy[tshift+i] = c
	}
	sort.Sort(cpuList(copy[tshift:]))
	if tshift != 0 {
		copy[0] = o.systemCPU.list[llen-tshift] // last, "cpu-total", becomes first
	}
	return copy
}

func (o *ostent) writeSystemCPU(cpuno int, m telegraf.Metric) {
	fields := m.Fields()
	id := m.Tags()["cpu"]
	c := cpu{N: id}

	for _, pair := range []struct {
		name  string
		value *int
	}{
		{"usage_user", &c.UserPct},
		{"usage_system", &c.SysPct},
		{"usage_iowait", &c.WaitPct},
		{"usage_idle", &c.IdlePct},
	} {
		if field, ok := fields[pair.name]; ok {
			if v, ok := field.(float64); ok {
				*pair.value = int(v)
			}
		}
	}

	o.systemCPU.mutex.Lock()
	defer o.systemCPU.mutex.Unlock()
	if len(o.systemCPU.list) < cpuno {
		list := make([]cpu, cpuno)
		copy(list, o.systemCPU.list)
		o.systemCPU.list = list
	}
	o.systemCPU.list[cpuno-1] = c
}

func (o *ostent) setCPUno(cpuno int) {
	o.systemCPU.mutex.Lock()
	defer o.systemCPU.mutex.Unlock()
	if len(o.systemCPU.list) > cpuno {
		o.systemCPU.list = o.systemCPU.list[:cpuno]
	}
}

func (o *ostent) writeSystemOstent(m telegraf.Metric) {
	o.systemOstent.mutex.Lock()
	defer o.systemOstent.mutex.Unlock()
	for k, field := range m.Fields() {
		if v, ok := field.(string); ok {
			o.systemOstent.kv[k] = v
		}
	}
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

	cpus := 0
	for _, m := range ms {
		switch m.Name() {
		case "system_ostent":
			o.writeSystemOstent(m)
		case "cpu":
			cpus++
			o.writeSystemCPU(cpus, m)
		}
		// default: if err := o.writeMetric(m); err != nil { return err }
	}
	o.setCPUno(cpus)
	return nil
}

var Output = &ostent{
	// Metrics: make(map[string]*Metric),
	systemOstent: group{kv: make(map[string]string)},
	systemCPU:    groupCPU{},
}

func init() { outputs.Add("ostent", func() telegraf.Output { return Output }) }
