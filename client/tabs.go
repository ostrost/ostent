package client

import (
	"github.com/ostrost/ostent/types"
)

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
	DFinodes types.SEQ
	DFbytes  types.SEQ

	DFinodesTitle string
	DFbytesTitle  string
}

// Title returns a label. "" return denotes unexpected error.
func (df DFtabs) Title(s types.SEQ) string {
	switch {
	case s == df.DFinodes:
		return df.DFinodesTitle
	case s == df.DFbytes:
		return df.DFbytesTitle
	}
	return ""
}

type IFtabs struct {
	IFpackets types.SEQ
	IFerrors  types.SEQ
	IFbytes   types.SEQ

	IFpacketsTitle string
	IFerrorsTitle  string
	IFbytesTitle   string
}

// Title returns a label. "" return denotes unexpected error.
func (fi IFtabs) Title(s types.SEQ) string {
	switch {
	case s == fi.IFpackets:
		return fi.IFpacketsTitle
	case s == fi.IFerrors:
		return fi.IFerrorsTitle
	case s == fi.IFbytes:
		return fi.IFbytesTitle
	}
	return ""
}

const (
	____IFTABID types.SEQ = iota
	IFPACKETS_TABID
	IFERRORS_TABID
	IFBYTES_TABID
)

const (
	____DFTABID types.SEQ = iota
	DFINODES_TABID
	DFBYTES_TABID
)

/* UNUSED ?
var IF_TABS = []types.SEQ{
	IFPACKETS_TABID,
	 IFERRORS_TABID,
	  IFBYTES_TABID,
}

var DF_TABS = []types.SEQ{
	DFINODES_TABID,
	 DFBYTES_TABID,
}
*/

var DFBIMAP = types.Seq2bimap(DFFS, // the default seq for ordering
	types.Seq2string{
		DFFS:    "fs",
		DFSIZE:  "size",
		DFUSED:  "used",
		DFAVAIL: "avail",
		DFMP:    "mp",
	}, []types.SEQ{
		DFFS, DFMP,
	})

var PSBIMAP = types.Seq2bimap(PSPID, // the default seq for ordering
	types.Seq2string{
		PSPID:  "pid",
		PSPRI:  "pri",
		PSNICE: "nice",
		PSSIZE: "size",
		PSRES:  "res",
		PSTIME: "time",
		PSNAME: "name",
		PSUID:  "user",
	}, []types.SEQ{
		PSNAME, PSUID,
	})

const (
	____DFIOTA types.SEQ = iota
	DFFS
	DFSIZE
	DFUSED
	DFAVAIL
	DFMP
)

const (
	____PSIOTA types.SEQ = iota
	PSPID
	PSPRI
	PSNICE
	PSSIZE
	PSRES
	PSTIME
	PSNAME
	PSUID
)
