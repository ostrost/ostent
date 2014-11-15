// +build linux

package cpu

import (
	"github.com/ostrost/ostent/types"
	sigar "github.com/rzab/gosigar"
)

func CalcTotal(cpu sigar.Cpu) uint64 {
	return cpu.Total() // gosigar implementation aka:
	// 	return cpu.User + cpu.Nice + cpu.Sys + cpu.Idle +
	// 		cpu.Wait + cpu.Irq + cpu.SoftIrq + cpu.Stolen
}

func (se Send) Fields() []types.NameString {
	cpu := se.raw()
	return []types.NameString{
		{"user", se.fraction(cpu.User)},
		{"nice", se.fraction(cpu.Nice)},
		{"system", se.fraction(cpu.Sys)},
		{"idle", se.fraction(cpu.Idle)},

		{"wait", se.fraction(cpu.Wait)},
		{"interrupt", se.fraction(cpu.Irq)},
		{"softirq", se.fraction(cpu.SoftIrq)},
		{"steal", se.fraction(cpu.Stolen)},
	}
}
