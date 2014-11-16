package cpu

// CPUInfo type has a list of CoreInfo.
type CPUInfo struct {
	List []CoreInfo // TODO rename to Cores
}

// CoreInfo type is a struct of core metrics.
type CoreInfo struct {
	N         string
	User      uint // percent without "%"
	Sys       uint // percent without "%"
	Idle      uint // percent without "%"
	UserClass string
	SysClass  string
	IdleClass string
	// UserSpark string
	// SysSpark  string
	// IdleSpark string
}
