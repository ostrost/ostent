// +build darwin

package cpu

import (
	"github.com/ostrost/ostent/types"
	sigar "github.com/rzab/gosigar"
)

func CalcTotal(cpu sigar.Cpu) uint64 {
	return cpu.User + cpu.Nice + cpu.Sys + cpu.Idle
	// gosigar cpu.Total() implementation adds .{Wait,{,Soft}Irq,Stolen} which is zero for darwin
}

func (se Send) Fields() []types.NameString {
	cpu := se.raw()
	return []types.NameString{
		{"user", se.fraction(cpu.User)},
		{"nice", se.fraction(cpu.Nice)},
		{"system", se.fraction(cpu.Sys)},
		{"idle", se.fraction(cpu.Idle)},
	}
}
