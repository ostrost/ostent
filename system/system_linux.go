// +build linux

package system

import (
	sigar "github.com/ostrost/gosigar"
	metrics "github.com/rcrowley/go-metrics"
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

func (emr *ExtraMetricRAM) UpdateRAM(got sigar.Mem, extra1, extra2 uint64) {
	emr.Used.Update(int64(got.ActualUsed))
	emr.Buffered.Update(int64(extra1))
	emr.Cached.Update(int64(extra2))
}

/* **************************************************************** */

type ExtraMetricCPU struct {
	Wait    *GaugePercent
	Irq     *GaugePercent
	SoftIrq *GaugePercent
	Stolen  *GaugePercent
}

func (emc *ExtraMetricCPU) UpdateCPU(sigarCpu sigar.Cpu, totalDelta int64) {
	emc.Wait.UpdatePercent(totalDelta, sigarCpu.Wait)
	emc.Irq.UpdatePercent(totalDelta, sigarCpu.Irq)
	emc.SoftIrq.UpdatePercent(totalDelta, sigarCpu.SoftIrq)
	emc.Stolen.UpdatePercent(totalDelta, sigarCpu.Stolen)
}

func NewMetricCPU(r metrics.Registry, name string) *MetricCPU {
	return ExtraNewMetricCPU(r, name, &ExtraMetricCPU{
		Wait:    NewGaugePercent(name+".wait", r),
		Irq:     NewGaugePercent(name+".interrupt", r),
		SoftIrq: NewGaugePercent(name+".softirq", r),
		Stolen:  NewGaugePercent(name+".steal", r),
	})
}
