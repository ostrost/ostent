// Package client is all about client state.
package client

import "github.com/ostrost/ostent/client/enums"

// Tab shadows enums.Uint and has Title string.
type Tab struct {
	enums.Uint
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
	DFINODES enums.Uint = iota
	DFBYTES
)

// Constants for IF tabs.
const (
	IFPACKETS enums.Uint = iota
	IFERRORS
	IFBYTES
)
