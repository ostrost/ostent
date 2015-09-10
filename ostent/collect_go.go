// +build !cgo

package ostent

import (
	"sync"
)

// IF is no-op without cgo.
func (m Machine) IF(_ Registry, wg *sync.WaitGroup) { wg.Done() }
