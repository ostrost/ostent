package cmdcobra

import (
	"sync"

	"github.com/spf13/cobra"
)

var (
	// PostRuns keeps a list of funcs to be RunE with cobra.Command's PostRunE.
	PostRuns Runs
	// PreRuns keeps a list of funcs to be RunE with cobra.Command's PreRunE.
	PreRuns Runs
)

// Runs keeps a list of funcs.
type Runs struct {
	Mutex sync.Mutex
	List  []func() error
}

// Add appends f to rs.List.
func (rs *Runs) Add(f func() error) {
	rs.Mutex.Lock()
	defer rs.Mutex.Unlock()
	rs.List = append(rs.List, f)
}

// Adds appends fs... to rs.List.
func (rs *Runs) Adds(fs ...func() error) {
	rs.Mutex.Lock()
	defer rs.Mutex.Unlock()
	rs.List = append(rs.List, fs...)
}

// RunE runs rs.List entries.
func (rs *Runs) RunE(*cobra.Command, []string) error {
	rs.Mutex.Lock()
	defer rs.Mutex.Unlock()
	for _, run := range rs.List {
		if err := run(); err != nil {
			return err
		}
	}
	return nil
}
