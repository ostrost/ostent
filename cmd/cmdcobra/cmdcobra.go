package cmdcobra

import (
	"sync"

	"github.com/spf13/cobra"
)

var (
	// PersistentPostRuns keeps a list of funcs to be cobra.Command's PersistentPostRunE.
	PersistentPostRuns Runs
	// PersistentPreRuns keeps a list of funcs to be cobra.Command's PersistentPreRunEE.
	PersistentPreRuns Runs
	// PreRuns keeps a list of funcs to be cobra.Command's PreRunE.
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
