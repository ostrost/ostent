package internal_models

import (
	"github.com/influxdata/telegraf"

	"github.com/ostrost/ostent/internal/buffer"
)

type RunningOutput struct {
	Output telegraf.Output
	Name   string
	// Config struct { Interval time.Duration }
	BufBatchSize int
	Buf          *buffer.Buffer
}

func (ro *RunningOutput) AddMetric(m telegraf.Metric) {
	ro.Buf.Add(m)
	if ro.Buf.Len() == ro.BufBatchSize {
		b := ro.Buf.Batch(ro.BufBatchSize)
		if err := ro.write(b); err != nil {
			panic(err)
		}
	}
}

func (ro *RunningOutput) Write() error {
	return ro.write(ro.Buf.Batch(ro.BufBatchSize))
}

func (ro *RunningOutput) write(ms []telegraf.Metric) error {
	if ms == nil || len(ms) == 0 {
		return nil
	}
	return ro.Output.Write(ms)
}
