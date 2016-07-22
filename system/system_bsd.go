// +build freebsd darwin

package system

import (
	metrics "github.com/rcrowley/go-metrics"
	"github.com/shirou/gopsutil/mem"
)

type ExtraMetricRAM struct {
	Inactive metrics.Gauge
	Wired    metrics.Gauge
	Active   metrics.Gauge
}

func NewMetricRAM(r metrics.Registry) *MetricRAM {
	return ExtraNewMetricRAM(r, &ExtraMetricRAM{
		Inactive: metrics.NewRegisteredGauge("memory.memory-inactive", r),
		Wired:    metrics.NewRegisteredGauge("memory.memory-wired", r),
		Active:   metrics.NewRegisteredGauge("memory.memory-active", r),
	})
}

func (emr *ExtraMetricRAM) UpdateRAM(stat *mem.VirtualMemoryStat) {
	emr.Inactive.Update(int64(stat.Inactive))
	emr.Wired.Update(int64(stat.Wired))
	emr.Active.Update(int64(stat.Active))
}

func NewMetricCPU(r metrics.Registry, name string) *MetricCPU {
	return ExtraNewMetricCPU(r, name, nil)
}
