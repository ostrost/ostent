// +build linux

package types

import (
	metrics "github.com/rcrowley/go-metrics"
	sigar "github.com/rzab/gosigar"
)

func CPUTotal(cpu sigar.Cpu) uint64 {
	return cpu.Total() // gosigar implementation aka:
	// 	return cpu.User + cpu.Nice + cpu.Sys + cpu.Idle +
	// 		cpu.Wait + cpu.Irq + cpu.SoftIrq + cpu.Stolen
}

type GaugeRAM struct {
	GaugeRAMCommon
	Free     metrics.Gauge
	Used     metrics.Gauge
	Buffered metrics.Gauge
	Cached   metrics.Gauge
}

func NewGaugeRAM(r metrics.Registry) GaugeRAM {
	return GaugeRAM{
		GaugeRAMCommon: NewGaugeRAMCommon(),
		Free:           metrics.NewRegisteredGauge("memory.memory-free", r),
		Used:           metrics.NewRegisteredGauge("memory.memory-used", r),
		Buffered:       metrics.NewRegisteredGauge("memory.memory-buffered", r),
		Cached:         metrics.NewRegisteredGauge("memory.memory-cached", r),
	}
}

func (gr *GaugeRAM) Update(got sigar.Mem, extra1, extra2 uint64) {
	gr.GaugeRAMCommon.UpdateCommon(got)
	gr.Free.Update(int64(got.Free))
	gr.Used.Update(int64(got.ActualUsed))
	gr.Buffered.Update(int64(extra1))
	gr.Cached.Update(int64(extra2))
}

type MetricCPU struct {
	MetricCPUCommon
	Wait    GaugePercent
	Irq     GaugePercent
	SoftIrq GaugePercent
	Stolen  GaugePercent
}

func (mc *MetricCPU) Update(sigarCpu sigar.Cpu) {
	totalDelta := mc.UpdateCommon(sigarCpu)
	mc.Wait.UpdatePercent(totalDelta, sigarCpu.Wait)
	mc.Irq.UpdatePercent(totalDelta, sigarCpu.Irq)
	mc.SoftIrq.UpdatePercent(totalDelta, sigarCpu.SoftIrq)
	mc.Stolen.UpdatePercent(totalDelta, sigarCpu.Stolen)
}

func NewMetricCPU(r metrics.Registry, name string) MetricCPU {
	return MetricCPU{
		MetricCPUCommon: NewMetricCPUCommon(r, name),
		Wait:            NewGaugePercent(name+".wait", r),
		Irq:             NewGaugePercent(name+".interrupt", r),
		SoftIrq:         NewGaugePercent(name+".softirq", r),
		Stolen:          NewGaugePercent(name+".steal", r),
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
