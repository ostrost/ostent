// +build linux

package types

import (
	"fmt"

	metrics "github.com/rcrowley/go-metrics"
	sigar "github.com/rzab/gosigar"
)

func RAMFields(ram RAM) []NameString {
	return []NameString{
		{"memory-free", fmt.Sprintf("%d", ram.Raw.Free)},
		{"memory-used", fmt.Sprintf("%d", ram.Raw.ActualUsed)},
		{"memory-buffered", fmt.Sprintf("%d", ram.Extra1)},
		{"memory-cached", fmt.Sprintf("%d", ram.Extra2)},
	}
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
