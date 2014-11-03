// +build linux

package cpu

import "github.com/ostrost/ostent/types"

func (se Send) calcTotal() uint64 {
	return se.cpu.User + se.cpu.Nice + se.cpu.Sys + se.cpu.Idle
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
