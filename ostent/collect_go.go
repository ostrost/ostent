// +build !cgo

package ostent

import (
	"sync"
)

// Interfaces is no-op without cgo.
func (m Machine) Interfaces(_ Registry, _ S2SRegistry, wg *sync.WaitGroup) {
	wg.Done()
}
