// +build darwin

package cpu

import (
	"github.com/ostrost/ostent/types"
	sigar "github.com/rzab/gosigar"
)

func cpuTotal(c *sigar.Cpu) uint64 {
	return c.User + c.Nice + c.Sys + c.Idle
}

func cpuFields(tc totalCpu) []types.NameFloat64 {
	return []types.NameFloat64{
		{"user", tc.fraction(tc.User)},
		{"nice", tc.fraction(tc.Nice)},
		{"system", tc.fraction(tc.Sys)},
		{"idle", tc.fraction(tc.Idle)},
	}
}
