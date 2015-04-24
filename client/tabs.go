//go:generate sh -c "jsonenums -type=UintDF; jsonenums -type=UintPS"
package client

var DFTABS = DFtabs{
	DFinodes: DFINODES_TABID,
	DFbytes:  DFBYTES_TABID,

	DFinodesTitle: "Disks inodes",
	DFbytesTitle:  "Disks",
}

var IFTABS = IFtabs{
	IFpackets: IFPACKETS_TABID,
	IFerrors:  IFERRORS_TABID,
	IFbytes:   IFBYTES_TABID,

	IFpacketsTitle: "Interfaces packets",
	IFerrorsTitle:  "Interfaces errors",
	IFbytesTitle:   "Interfaces",
}

type DFtabs struct {
	DFinodes Uint
	DFbytes  Uint

	DFinodesTitle string
	DFbytesTitle  string
}

// Title returns a label. "" return denotes unidentified p.
func (df DFtabs) Title(u Uint) string {
	switch {
	case u == df.DFinodes:
		return df.DFinodesTitle
	case u == df.DFbytes:
		return df.DFbytesTitle
	}
	return ""
}

type IFtabs struct {
	IFpackets Uint
	IFerrors  Uint
	IFbytes   Uint

	IFpacketsTitle string
	IFerrorsTitle  string
	IFbytesTitle   string
}

// Title returns a label. "" return denotes unidentified p.
func (fi IFtabs) Title(u Uint) string {
	switch {
	case u == fi.IFpackets:
		return fi.IFpacketsTitle
	case u == fi.IFerrors:
		return fi.IFerrorsTitle
	case u == fi.IFbytes:
		return fi.IFbytesTitle
	}
	return ""
}

const (
	IFPACKETS_TABID Uint = iota
	IFERRORS_TABID
	IFBYTES_TABID
)

const (
	DFINODES_TABID Uint = iota
	DFBYTES_TABID
)

/* UNUSED ?
var IF_TABS = []Uint{
	IFPACKETS_TABID,
	 IFERRORS_TABID,
	  IFBYTES_TABID,
}

var DF_TABS = []Uint{
	DFINODES_TABID,
	 DFBYTES_TABID,
}
*/

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
)

// RenamedConstError denotes an error.
type RenamedConstError string

func (rc RenamedConstError) Error() string { return string(rc) }

// Uint-derived types:

// UintDF is a derived Uint for constants.
type UintDF Uint

// UintPS is a derived Uint for constants.
type UintPS Uint
