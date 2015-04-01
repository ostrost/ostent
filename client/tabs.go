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
	DFinodes SEQ
	DFbytes  SEQ

	DFinodesTitle string
	DFbytesTitle  string
}

// Title returns a label. "" return denotes unexpected error.
func (df DFtabs) Title(s SEQ) string {
	switch {
	case s == df.DFinodes:
		return df.DFinodesTitle
	case s == df.DFbytes:
		return df.DFbytesTitle
	}
	return ""
}

type IFtabs struct {
	IFpackets SEQ
	IFerrors  SEQ
	IFbytes   SEQ

	IFpacketsTitle string
	IFerrorsTitle  string
	IFbytesTitle   string
}

// Title returns a label. "" return denotes unexpected error.
func (fi IFtabs) Title(s SEQ) string {
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
	____IFTABID SEQ = iota
	IFPACKETS_TABID
	IFERRORS_TABID
	IFBYTES_TABID
)

const (
	____DFTABID SEQ = iota
	DFINODES_TABID
	DFBYTES_TABID
)

/* UNUSED ?
var IF_TABS = []SEQ{
	IFPACKETS_TABID,
	 IFERRORS_TABID,
	  IFBYTES_TABID,
}

var DF_TABS = []SEQ{
	DFINODES_TABID,
	 DFBYTES_TABID,
}
*/

var DFBIMAP = Seq2bimap(SEQ(DFFS), // the default seq for ordering
	Seq2string{
		SEQ(DFFS):    "fs",
		SEQ(DFMP):    "mp",
		SEQ(DFSIZE):  "size",
		SEQ(DFUSED):  "used",
		SEQ(DFAVAIL): "avail",
	}, []SEQ{
		SEQ(DFFS), SEQ(DFMP),
	})

var PSBIMAP = Seq2bimap(PSPID, // the default seq for ordering
	Seq2string{
		PSPID:  "pid",
		PSPRI:  "pri",
		PSNICE: "nice",
		PSSIZE: "size",
		PSRES:  "res",
		PSTIME: "time",
		PSNAME: "name",
		PSUID:  "user",
	}, []SEQ{
		PSNAME, PSUID,
	})

const (
	____DFIOTA SEQ = iota // TODO rename to DFZERO
	DFFS
	DFMP
	DFSIZE
	DFUSED
	DFAVAIL

	// DEFDFFS defines default DFSEQ.
	// The default is to be omitted from link parameters.
	DEFDFFS = DFFS
)

const (
	____PSIOTA SEQ = iota
	PSPID
	PSPRI
	PSNICE
	PSSIZE
	PSRES
	PSTIME
	PSNAME
	PSUID
)
