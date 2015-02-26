// +build freebsd darwin

package system

import (
	"github.com/ostrost/ostent/types"
	metrics "github.com/rcrowley/go-metrics"
	sigar "github.com/rzab/gosigar"
)

type MetricRAM struct {
	types.MetricRAMCommon
	Free     metrics.Gauge
	Inactive metrics.Gauge
	Wired    metrics.Gauge
	Active   metrics.Gauge
}

func NewMetricRAM(r metrics.Registry) MetricRAM {
	return MetricRAM{
		MetricRAMCommon: types.NewMetricRAMCommon(),
		Free:            metrics.NewRegisteredGauge("memory.memory-free", r),
		Inactive:        metrics.NewRegisteredGauge("memory.memory-inactive", r),
		Wired:           metrics.NewRegisteredGauge("memory.memory-wired", r),
		Active:          metrics.NewRegisteredGauge("memory.memory-active", r),
	}
}

func (mr *MetricRAM) Update(got sigar.Mem, extra1, extra2 uint64) {
	mr.MetricRAMCommon.UpdateCommon(got)
	mr.Free.Update(int64(got.Free))
	mr.Inactive.Update(int64(got.ActualFree - got.Free))
	mr.Wired.Update(int64(extra1))
	mr.Active.Update(int64(extra2))
}

type MetricCPU struct {
	*types.MetricCPUCommon
}

func (mc *MetricCPU) Update(cpu sigar.Cpu) {
	total := cpu.User + cpu.Nice + cpu.Sys + cpu.Idle
	// gosigar cpu.Total() implementation adds .{Wait,{,Soft}Irq,Stolen}
	// which is zero for darwin
	mc.UpdateCommon(cpu, total)
}

func NewMetricCPU(r metrics.Registry, name string) *MetricCPU {
	return &MetricCPU{
		MetricCPUCommon: types.NewMetricCPUCommon(r, name),
	}
}

func CPUAdd(sum *sigar.Cpu, other sigar.Cpu) {
	sum.User += other.User
	sum.Nice += other.Nice
	sum.Sys += other.Sys
	sum.Idle += other.Idle
}
