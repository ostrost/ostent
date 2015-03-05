// +build linux

package system

import (
	"github.com/ostrost/ostent/system/operating"
	metrics "github.com/rcrowley/go-metrics"
	sigar "github.com/rzab/gosigar"
)

type ExtraMetricRAM struct {
	Used     metrics.Gauge
	Buffered metrics.Gauge
	Cached   metrics.Gauge
}

func NewMetricRAM(r metrics.Registry) *operating.MetricRAM {
	return operating.ExtraNewMetricRAM(r, &ExtraMetricRAM{
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
	Wait    *operating.GaugePercent
	Irq     *operating.GaugePercent
	SoftIrq *operating.GaugePercent
	Stolen  *operating.GaugePercent
}

func (emc *ExtraMetricCPU) UpdateCPU(sigarCpu sigar.Cpu, totalDelta int64) {
	emc.Wait.UpdatePercent(totalDelta, sigarCpu.Wait)
	emc.Irq.UpdatePercent(totalDelta, sigarCpu.Irq)
	emc.SoftIrq.UpdatePercent(totalDelta, sigarCpu.SoftIrq)
	emc.Stolen.UpdatePercent(totalDelta, sigarCpu.Stolen)
}

func NewMetricCPU(r metrics.Registry, name string) *operating.MetricCPU {
	return operating.ExtraNewMetricCPU(r, name, &ExtraMetricCPU{
		Wait:    operating.NewGaugePercent(name+".wait", r),
		Irq:     operating.NewGaugePercent(name+".interrupt", r),
		SoftIrq: operating.NewGaugePercent(name+".softirq", r),
		Stolen:  operating.NewGaugePercent(name+".steal", r),
	})
}
