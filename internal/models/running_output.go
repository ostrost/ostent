package internal_models

import (
	"github.com/influxdata/telegraf"

	"github.com/ostrost/ostent/internal/buffer"
)

// RunningOutput contains the output configuration
type RunningOutput struct {
	Name            string
	Output          telegraf.Output
	MetricBatchSize int

	metrics *buffer.Buffer
}

func NewRunningOutput(
	name string,
	output telegraf.Output,
) *RunningOutput {
	batchSize := 1000
	ro := &RunningOutput{
		Name:            name,
		metrics:         buffer.NewBuffer(batchSize),
		Output:          output,
		MetricBatchSize: batchSize,
	}
	return ro
}

// AddMetric adds a metric to the output. This function can also write cached
// points if FlushBufferWhenFull is true.
func (ro *RunningOutput) AddMetric(metric telegraf.Metric) {
	ro.metrics.Add(metric)
	if ro.metrics.Len() == ro.MetricBatchSize {
		batch := ro.metrics.Batch(ro.MetricBatchSize)
		err := ro.write(batch)
		if err != nil {
			panic(err)
		}
	}
}

// Write writes all cached points to this output.
func (ro *RunningOutput) Write() error {
	var err error

	batch := ro.metrics.Batch(ro.MetricBatchSize)
	if err == nil {
		err = ro.write(batch)
	}
	if err != nil {
		return err
	}
	return nil
}

func (ro *RunningOutput) write(metrics []telegraf.Metric) error {
	if metrics == nil || len(metrics) == 0 {
		return nil
	}
	err := ro.Output.Write(metrics)
	return err
}
