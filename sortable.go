package ostent

import (
	"github.com/ostrost/ostent/types"
)

type interfaceOrder []types.Interface

func (io interfaceOrder) Len() int {
	return len(io)
}

func (io interfaceOrder) Swap(i, j int) {
	io[i], io[j] = io[j], io[i]
}

func (io interfaceOrder) Less(i, j int) bool {
	if rx_lo.Match([]byte(io[i].NameKey)) {
		return false
	}
	return io[i].NameKey < io[j].NameKey
}
