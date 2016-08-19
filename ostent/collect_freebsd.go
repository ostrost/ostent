// +build freebsd

package ostent

// ProcName returns procName back.
func ProcName(_ int, procName string) string { return procName }
