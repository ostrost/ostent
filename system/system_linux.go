// +build linux

package system

import (
	metrics "github.com/rcrowley/go-metrics"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

type ExtraMetricRAM struct {
	Used     metrics.Gauge
	Buffered metrics.Gauge
	Cached   metrics.Gauge
}

func NewMetricRAM(r metrics.Registry) *MetricRAM {
	return ExtraNewMetricRAM(r, &ExtraMetricRAM{
		Used:     metrics.NewRegisteredGauge("memory.memory-used", r),
		Buffered: metrics.NewRegisteredGauge("memory.memory-buffered", r),
		Cached:   metrics.NewRegisteredGauge("memory.memory-cached", r),
	})
}

func (emr *ExtraMetricRAM) UpdateRAM(stat *mem.VirtualMemoryStat) {
	emr.Used.Update(int64(stat.Used))
	emr.Buffered.Update(int64(stat.Buffers))
	emr.Cached.Update(int64(stat.Cached))
}

/* **************************************************************** */

type ExtraMetricCPU struct {
	Wait    *GaugePercent
	Irq     *GaugePercent
	SoftIrq *GaugePercent
	Stolen  *GaugePercent
}

func (emc *ExtraMetricCPU) UpdateCPU(stat cpu.TimesStat, totalDelta float64) {
	emc.Wait.UpdatePercent(totalDelta, stat.Iowait)
	emc.Irq.UpdatePercent(totalDelta, stat.Irq)
	emc.SoftIrq.UpdatePercent(totalDelta, stat.Softirq)
	emc.Stolen.UpdatePercent(totalDelta, stat.Stolen)
}

func NewMetricCPU(r metrics.Registry, name string) *MetricCPU {
	return ExtraNewMetricCPU(r, name, &ExtraMetricCPU{
		Wait:    NewGaugePercent(name+".wait", r),
		Irq:     NewGaugePercent(name+".interrupt", r),
		SoftIrq: NewGaugePercent(name+".softirq", r),
		Stolen:  NewGaugePercent(name+".steal", r),
	})
}
