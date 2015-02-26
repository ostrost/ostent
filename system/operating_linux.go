// +build linux

package system

import (
	"github.com/ostrost/ostent/types"
	metrics "github.com/rcrowley/go-metrics"
	sigar "github.com/rzab/gosigar"
)

type MetricRAM struct {
	types.MetricRAMCommon
	Free     metrics.Gauge
	Used     metrics.Gauge
	Buffered metrics.Gauge
	Cached   metrics.Gauge
}

func NewMetricRAM(r metrics.Registry) MetricRAM {
	return MetricRAM{
		MetricRAMCommon: types.NewMetricRAMCommon(),
		Free:            metrics.NewRegisteredGauge("memory.memory-free", r),
		Used:            metrics.NewRegisteredGauge("memory.memory-used", r),
		Buffered:        metrics.NewRegisteredGauge("memory.memory-buffered", r),
		Cached:          metrics.NewRegisteredGauge("memory.memory-cached", r),
	}
}

func (mr *MetricRAM) Update(got sigar.Mem, extra1, extra2 uint64) {
	mr.MetricRAMCommon.UpdateCommon(got)
	mr.Free.Update(int64(got.Free))
	mr.Used.Update(int64(got.ActualUsed))
	mr.Buffered.Update(int64(extra1))
	mr.Cached.Update(int64(extra2))
}

type MetricCPU struct {
	*types.MetricCPUCommon
	Wait    *types.GaugePercent
	Irq     *types.GaugePercent
	SoftIrq *types.GaugePercent
	Stolen  *types.GaugePercent
}

func (mc *MetricCPU) Update(sigarCpu sigar.Cpu) {
	total := sigarCpu.Total() // gosigar implementation aka:
	// .User + .Nice + .Sys + .Idle + .Wait + .Irq + .SoftIrq + .Stolen
	totalDelta := mc.UpdateCommon(sigarCpu, total)
	mc.Wait.UpdatePercent(totalDelta, sigarCpu.Wait)
	mc.Irq.UpdatePercent(totalDelta, sigarCpu.Irq)
	mc.SoftIrq.UpdatePercent(totalDelta, sigarCpu.SoftIrq)
	mc.Stolen.UpdatePercent(totalDelta, sigarCpu.Stolen)
}

func NewMetricCPU(r metrics.Registry, name string) *MetricCPU {
	return &MetricCPU{
		MetricCPUCommon: types.NewMetricCPUCommon(r, name),
		Wait:            types.NewGaugePercent(name+".wait", r),
		Irq:             types.NewGaugePercent(name+".interrupt", r),
		SoftIrq:         types.NewGaugePercent(name+".softirq", r),
		Stolen:          types.NewGaugePercent(name+".steal", r),
	}
}

func CPUAdd(sum *sigar.Cpu, other sigar.Cpu) {
	sum.User += other.User
	sum.Nice += other.Nice
	sum.Sys += other.Sys
	sum.Idle += other.Idle
	sum.Wait += other.Wait
	sum.Irq += other.Irq
	sum.SoftIrq += other.SoftIrq
	sum.Stolen += other.Stolen
}
