package registry

import (
	"github.com/ostrost/ostent/getifaddrs"
	sigar "github.com/rzab/gosigar"
)

type Registry interface {
	UpdateIFdata(getifaddrs.IfData)
	UpdateCPU([]sigar.Cpu)
	UpdateLoadAverage(sigar.LoadAverage)
	UpdateSwap(sigar.Swap)
	UpdateRAM(sigar.Mem, uint64, uint64)
	UpdateDF(sigar.FileSystem, sigar.FileSystemUsage)
}
