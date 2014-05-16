package ostential
import (
	"ostential/types"
)

type internalClient struct {
	// NB lowercase fields only, NOT to be marshalled/exported

	psNotexpandable *bool
	psLimit int
}

type clientState struct {
	internalClient `json:"-"` // NB not marshalled

	HideMEM *bool `json:",omitempty"`
	HideIF  *bool `json:",omitempty"`
	HideCPU *bool `json:",omitempty"`
	HideDF  *bool `json:",omitempty"`
	HidePS  *bool `json:",omitempty"`

	HideSWAP  *bool `json:",omitempty"`
	ExpandIF  *bool `json:",omitempty"`
	ExpandCPU *bool `json:",omitempty"`
	ExpandDF  *bool `json:",omitempty"`

	TabIF *types.SEQ `json:",omitempty"`
	TabDF *types.SEQ `json:",omitempty"`

	IFTABS *iftabs `json:",omitempty"` // immutable, constant
	DFTABS *dftabs `json:",omitempty"` // immutable, constant

	// PSusers []string `json:omitempty`

	ConfigMEM *bool `json:",omitempty"`
	ConfigIF  *bool `json:",omitempty"`
	ConfigCPU *bool `json:",omitempty"`
	ConfigDF  *bool `json:",omitempty"`
	ConfigPS  *bool `json:",omitempty"`
}

type dftabs struct {
	DFbytes  types.SEQ
	DFinodes types.SEQ
}

type iftabs struct {
	IFpackets types.SEQ
	IFerrors  types.SEQ
	IFbytes   types.SEQ
}

func(_  clientState) merge_bool(dest, src *bool)    { if src != nil { *dest = *src } }
func(_  clientState) mergeSEQ(dest, src *types.SEQ) { if src != nil { *dest = *src } }

func(cs *clientState) Merge(ps clientState) {
	cs.merge_bool(cs.HideMEM, ps.HideMEM)
	cs.merge_bool(cs.HideIF,  ps.HideIF)
	cs.merge_bool(cs.HideCPU, ps.HideCPU)
	cs.merge_bool(cs.HideDF,  ps.HideDF)
	cs.merge_bool(cs.HidePS,  ps.HidePS)

	cs.merge_bool(cs.HideSWAP,  ps.HideSWAP)
	cs.merge_bool(cs.ExpandIF,  ps.ExpandIF)
	cs.merge_bool(cs.ExpandCPU, ps.ExpandCPU)
	cs.merge_bool(cs.ExpandDF,  ps.ExpandDF)

	cs.mergeSEQ(cs.TabIF, ps.TabIF)
	cs.mergeSEQ(cs.TabDF, ps.TabDF)

	cs.merge_bool(cs.ConfigMEM, ps.ConfigMEM)
	cs.merge_bool(cs.ConfigIF,  ps.ConfigIF)
	cs.merge_bool(cs.ConfigCPU, ps.ConfigCPU)
	cs.merge_bool(cs.ConfigDF,  ps.ConfigDF)
	cs.merge_bool(cs.ConfigPS,  ps.ConfigPS)
}

const (
	____IFTABID types.SEQ = iota
	IFPACKETS_TABID
	 IFERRORS_TABID
	  IFBYTES_TABID
)

var IF_TABS = []types.SEQ{
	IFPACKETS_TABID,
	 IFERRORS_TABID,
	  IFBYTES_TABID,
}

const (
	____DFTABID types.SEQ = iota
	DFINODES_TABID
	 DFBYTES_TABID
)

var DF_TABS = []types.SEQ{
	DFINODES_TABID,
	 DFBYTES_TABID,
}

func newfalse() *bool { return new(bool) }
func newtrue()  *bool { return newbool(true); }
func newbool(v bool) (b *bool) { b = new(bool); *b = v; return }

func newseq(v types.SEQ) *types.SEQ {
	s := new(types.SEQ)
	*s = v
	return s
}

func defaultClientState() clientState {
	cs := clientState{}

	cs.HideMEM = newfalse()
	cs.HideIF  = newfalse()
	cs.HideCPU = newfalse()
	cs.HideDF  = newfalse()
	cs.HidePS  = newfalse()

	cs.HideSWAP  = newfalse()
	cs.ExpandIF  = newfalse()
	cs.ExpandCPU = newfalse()
	cs.ExpandDF  = newfalse()

	cs.TabIF = newseq(IFBYTES_TABID)
	cs.TabDF = newseq(DFBYTES_TABID)

	cs.DFTABS = &dftabs{ // immutable, constant
		DFbytes:  DFBYTES_TABID,
		DFinodes: DFINODES_TABID,
	}

	cs.IFTABS = &iftabs{ // immutable, constant
		IFpackets: IFPACKETS_TABID,
		IFerrors:  IFERRORS_TABID,
		IFbytes:   IFBYTES_TABID,
	}

	// reverse logic
	configdefault := true
	// configdefault = false // DEVELOPMENT

	cs.ConfigMEM = newbool(configdefault)
	cs.ConfigIF  = newbool(configdefault)
	cs.ConfigCPU = newbool(configdefault)
	cs.ConfigDF  = newbool(configdefault)
	cs.ConfigPS  = newbool(configdefault)

	cs.psLimit = 16
	// cs.psNotexpandable = newfalse()

	return cs
}
