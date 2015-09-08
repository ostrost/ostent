// +build !cgo

package ostent

import (
	"sync"
)

// IF is no-op without cgo.
func (m Machine) IF(_ Registry, _ S2SRegistry, wg *sync.WaitGroup) {
	wg.Done()
}
