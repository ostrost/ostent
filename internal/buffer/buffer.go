package buffer

import "github.com/influxdata/telegraf"

type Buffer struct{ buf chan telegraf.Metric }

func NewBuffer(size int) *Buffer { return &Buffer{buf: make(chan telegraf.Metric, size)} }

func (b *Buffer) Len() int { return len(b.buf) }

func (b *Buffer) Add(ms ...telegraf.Metric) {
	for i := range ms {
		select {
		case b.buf <- ms[i]:
		default:
			<-b.buf
			b.buf <- ms[i]
		}
	}
}

func (b *Buffer) Batch(bsize int) []telegraf.Metric {
	n := b.Len()
	if n >= bsize {
		n = bsize
	}
	o := make([]telegraf.Metric, n)
	for i := 0; i < n; i++ {
		o[i] = <-b.buf
	}
	return o
}
