// +build darwin

package types

import (
	"fmt"

	metrics "github.com/rcrowley/go-metrics"
	sigar "github.com/rzab/gosigar"
)

func RAMFields(ram RAM) []NameString {
	return []NameString{
		{"memory-free", fmt.Sprintf("%d", ram.Raw.Free)},
		{"memory-inactive", fmt.Sprintf("%d", ram.Raw.ActualFree-ram.Raw.Free)},
		{"memory-wired", fmt.Sprintf("%d", ram.Extra1)},
		{"memory-active", fmt.Sprintf("%d", ram.Extra2)},
	}
}

type GaugeRAM struct {
	GaugeRAMCommon
	Free     metrics.Gauge
	Inactive metrics.Gauge
	Wired    metrics.Gauge
	Active   metrics.Gauge
}

func NewGaugeRAM(r metrics.Registry) GaugeRAM {
	return GaugeRAM{
		GaugeRAMCommon: NewGaugeRAMCommon(),
		Free:           metrics.NewRegisteredGauge("memory.memory-free", r),
		Inactive:       metrics.NewRegisteredGauge("memory.memory-inactive", r),
		Wired:          metrics.NewRegisteredGauge("memory.memory-wired", r),
		Active:         metrics.NewRegisteredGauge("memory.memory-active", r),
	}
}

func (gr *GaugeRAM) Update(got sigar.Mem, extra1, extra2 uint64) {
	gr.GaugeRAMCommon.UpdateCommon(got)
	gr.Free.Update(int64(got.Free))
	gr.Inactive.Update(int64(got.ActualFree - got.Free))
	gr.Wired.Update(int64(extra1))
	gr.Active.Update(int64(extra2))
}
