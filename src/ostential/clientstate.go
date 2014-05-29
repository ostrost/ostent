package ostential
import (
	"ostential/types"
)

type internalClient struct {
	// NB lowercase fields only, NOT to be marshalled/exported

	psNotexpandable *bool
	psLimit int
}

type title string
func (ti *title) merge_title(ns string, dt **title) {
	if string(*ti) == ns {
		return
	}
	*dt = (*title)(newstring(ns))
	*ti = **dt
}

type clientState struct {
	internalClient `json:"-"` // NB not marshalled

	HideMEM *bool `json:",omitempty"`
	HideIF  *bool `json:",omitempty"`
	HideCPU *bool `json:",omitempty"`
	HideDF  *bool `json:",omitempty"`
	HidePS  *bool `json:",omitempty"`
	HideVG  *bool `json:",omitempty"`

	HideSWAP  *bool `json:",omitempty"`
	ExpandIF  *bool `json:",omitempty"`
	ExpandCPU *bool `json:",omitempty"`
	ExpandDF  *bool `json:",omitempty"`

	TabIF *types.SEQ `json:",omitempty"`
	TabDF *types.SEQ `json:",omitempty"`
	TabIFtitle *title `json:",omitempty"`
	TabDFtitle *title `json:",omitempty"`

	// PSusers []string `json:omitempty`

	HideconfigMEM *bool `json:",omitempty"`
	HideconfigIF  *bool `json:",omitempty"`
	HideconfigCPU *bool `json:",omitempty"`
	HideconfigDF  *bool `json:",omitempty"`
	HideconfigPS  *bool `json:",omitempty"`
	HideconfigVG  *bool `json:",omitempty"`
}

func(_  clientState) merge_bool(dest, src *bool)    { if src != nil { *dest = *src } }
func(_  clientState) mergeSEQ(dest, src *types.SEQ) { if src != nil { *dest = *src } }

func(cs *clientState) Merge(ps clientState, ds *clientState) {
	cs.merge_bool(cs.HideMEM, ps.HideMEM)
	cs.merge_bool(cs.HideIF,  ps.HideIF)
	cs.merge_bool(cs.HideCPU, ps.HideCPU)
	cs.merge_bool(cs.HideDF,  ps.HideDF)
	cs.merge_bool(cs.HidePS,  ps.HidePS)
	cs.merge_bool(cs.HideVG,  ps.HideVG)

	cs.merge_bool(cs.HideSWAP,  ps.HideSWAP)
	cs.merge_bool(cs.ExpandIF,  ps.ExpandIF)
	cs.merge_bool(cs.ExpandCPU, ps.ExpandCPU)
	cs.merge_bool(cs.ExpandDF,  ps.ExpandDF)

	cs.mergeSEQ(cs.TabIF, ps.TabIF)
	cs.mergeSEQ(cs.TabDF, ps.TabDF)

	cs.TabIFtitle.merge_title(IFTABS.Title(*cs.TabIF), &ds.TabIFtitle)
	cs.TabDFtitle.merge_title(DFTABS.Title(*cs.TabDF), &ds.TabDFtitle)

	cs.merge_bool(cs.HideconfigMEM, ps.HideconfigMEM)
	cs.merge_bool(cs.HideconfigIF,  ps.HideconfigIF)
	cs.merge_bool(cs.HideconfigCPU, ps.HideconfigCPU)
	cs.merge_bool(cs.HideconfigDF,  ps.HideconfigDF)
	cs.merge_bool(cs.HideconfigPS,  ps.HideconfigPS)
	cs.merge_bool(cs.HideconfigVG,  ps.HideconfigVG)
}

func newstring(s string) *string {
	p := new(string)
	*p = s
	return p
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
	cs.HideVG  = newfalse()

	cs.HideSWAP  = newfalse()
	cs.ExpandIF  = newfalse()
	cs.ExpandCPU = newfalse()
	cs.ExpandDF  = newfalse()

	cs.TabIF = newseq(IFBYTES_TABID)
	cs.TabDF = newseq(DFBYTES_TABID)
	cs.TabIFtitle = (*title)(newstring(IFTABS.Title(*cs.TabIF)))
	cs.TabDFtitle = (*title)(newstring(DFTABS.Title(*cs.TabDF)))

	hideconfig := true
	// hideconfig = false // DEVELOPMENT

	cs.HideconfigMEM = newbool(hideconfig)
	cs.HideconfigIF  = newbool(hideconfig)
	cs.HideconfigCPU = newbool(hideconfig)
	cs.HideconfigDF  = newbool(hideconfig)
	cs.HideconfigPS  = newbool(hideconfig)
	cs.HideconfigVG  = newbool(hideconfig)

	cs.psLimit = 8
	// cs.psNotexpandable = newfalse()

	return cs
}
