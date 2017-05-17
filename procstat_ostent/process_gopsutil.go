// +build freebsd darwin

package procstat_ostent

func NewProc(pid PID) (Process, error) { return psutilNewProc(pid) }
