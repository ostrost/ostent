// +build linux

package system

import (
	"github.com/ostrost/ostent/types"
	metrics "github.com/rcrowley/go-metrics"
	sigar "github.com/rzab/gosigar"
)

type ExtraMetricRAM struct {
	Used     metrics.Gauge
	Buffered metrics.Gauge
	Cached   metrics.Gauge
}

func NewMetricRAM(r metrics.Registry) *types.MetricRAM {
	return types.ExtraNewMetricRAM(r, &ExtraMetricRAM{
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
	Wait    *types.GaugePercent
	Irq     *types.GaugePercent
	SoftIrq *types.GaugePercent
	Stolen  *types.GaugePercent
}

func (emc *ExtraMetricCPU) UpdateCPU(sigarCpu sigar.Cpu, totalDelta int64) {
	emc.Wait.UpdatePercent(totalDelta, sigarCpu.Wait)
	emc.Irq.UpdatePercent(totalDelta, sigarCpu.Irq)
	emc.SoftIrq.UpdatePercent(totalDelta, sigarCpu.SoftIrq)
	emc.Stolen.UpdatePercent(totalDelta, sigarCpu.Stolen)
}

func NewMetricCPU(r metrics.Registry, name string) *types.MetricCPU {
	return types.ExtraNewMetricCPU(r, name, &ExtraMetricCPU{
		Wait:    types.NewGaugePercent(name+".wait", r),
		Irq:     types.NewGaugePercent(name+".interrupt", r),
		SoftIrq: types.NewGaugePercent(name+".softirq", r),
		Stolen:  types.NewGaugePercent(name+".steal", r),
	})
}
