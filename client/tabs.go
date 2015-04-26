//go:generate sh -c "jsonenums -type=UintDF; jsonenums -type=UintPS"

// Package client is all about client state.
package client

// Tab shadows Uint and has Title string.
type Tab struct {
	Uint
	Title string
}

// Tabs is to define known tabs by string. In use in templates.
type Tabs map[string]Tab

// DFTABS is a map containing defined DF Tab's.
var DFTABS = Tabs{
	"dFINODES": {DFINODES, "Disks inodes"},
	"dFBYTES":  {DFBYTES, "Disks"},
}

// IFTABS is a map containing defined IF Tab's.
var IFTABS = Tabs{
	"iFPACKETS": {IFPACKETS, "Interfaces packets"},
	"iFERRORS":  {IFERRORS, "Interfaces errors"},
	"iFBYTES":   {IFBYTES, "Interfaces"},
}

// Constants for DF tabs.
const (
	DFINODES Uint = iota
	DFBYTES
)

// Constants for IF tabs.
const (
	IFPACKETS Uint = iota
	IFERRORS
	IFBYTES
)

// Constants for DF sorting criterion.
const (
	FS UintDF = iota
	MP
	TOTAL
	USED
	AVAIL
)

// Constants for PS sorting criterion.
const (
	PID UintPS = iota
	PRI
	NICE
	VIRT
	RES
	TIME
	NAME
	UID
	USER
)

// RenamedConstError denotes an error.
type RenamedConstError string

func (rc RenamedConstError) Error() string { return string(rc) }

// Uint-derived types:

// UintDF is a derived Uint for constants.
type UintDF Uint

// UintPS is a derived Uint for constants.
type UintPS Uint
