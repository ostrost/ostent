package system

//"github.com/ostrost/ostent/types"
//metrics "github.com/rcrowley/go-metrics"
//sigar "github.com/rzab/gosigar"

func (mr *MetricRAM) UsedValue() uint64 { // Total - Free
	return uint64(mr.Total.Snapshot().Value() - mr.Free.Snapshot().Value())
}
