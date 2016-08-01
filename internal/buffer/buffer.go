package buffer

import "github.com/influxdata/telegraf"

// Buffer is an object for storing metrics in a circular buffer.
type Buffer struct {
	buf chan telegraf.Metric
}

// NewBuffer returns a Buffer
//   size is the maximum number of metrics that Buffer will cache. If Add is
//   called when the buffer is full, then the oldest metric(s) will be dropped.
func NewBuffer(size int) *Buffer {
	return &Buffer{
		buf: make(chan telegraf.Metric, size),
	}
}

// Len returns the current length of the buffer.
func (b *Buffer) Len() int {
	return len(b.buf)
}

func (b *Buffer) Add(metrics ...telegraf.Metric) {
	for i, _ := range metrics {
		select {
		case b.buf <- metrics[i]:
		default:
			<-b.buf
			b.buf <- metrics[i]
		}
	}
}

// Batch returns a batch of metrics of size batchSize.
// the batch will be of maximum length batchSize. It can be less than batchSize,
// if the length of Buffer is less than batchSize.
func (b *Buffer) Batch(batchSize int) []telegraf.Metric {
	n := min(len(b.buf), batchSize)
	out := make([]telegraf.Metric, n)
	for i := 0; i < n; i++ {
		out[i] = <-b.buf
	}
	return out
}

func min(a, b int) int {
	if b < a {
		return b
	}
	return a
}
