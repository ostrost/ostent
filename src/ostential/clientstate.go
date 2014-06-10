package ostential
import (
	"ostential/types"

	"time"
	"strings"
	"encoding/json"
)

type refresh struct {
	time.Duration
	tick int
}
func(r refresh) MarshalJSON() ([]byte, error) {
	s := r.Duration.String()
	if strings.HasSuffix(s, "m0s") {
		s = strings.TrimSuffix(s, "0s")
	}
	if strings.HasSuffix(s, "h0m") {
		s = strings.TrimSuffix(s, "0m")
	}
	return json.Marshal(s)
}

func(r *refresh) refresh(forcerefresh bool) Boole {
	if forcerefresh {
		return Boole(true)
	}
	r.tick++
	if r.tick < int(r.Duration / time.Second) {
		return Boole(false)
	}
	r.tick = 0
	return Boole(true)
}

func(r refresh) expires() bool {
	return r.tick + 1 >= int(r.Duration / time.Second)
}

type internalClient struct {
	// NB lowercase fields only, NOT to be marshalled/exported

	psLimit int

	psSEQ types.SEQ
	dfSEQ types.SEQ
}

type title string
func (ti *title) merge(ns string, dt **title) {
	*dt = nil
	if string(*ti) == ns {
		return
	}
	*dt = newtitle(ns)
	*ti = **dt
}

type commonClient struct {
	HideMEM *Boole `json:",omitempty"`
	HideIF  *Boole `json:",omitempty"`
	HideCPU *Boole `json:",omitempty"`
	HideDF  *Boole `json:",omitempty"`
	HidePS  *Boole `json:",omitempty"`
	HideVG  *Boole `json:",omitempty"`

	HideSWAP  *Boole `json:",omitempty"`
	ExpandIF  *Boole `json:",omitempty"`
	ExpandCPU *Boole `json:",omitempty"`
	ExpandDF  *Boole `json:",omitempty"`

	TabIF *types.SEQ `json:",omitempty"`
	TabDF *types.SEQ `json:",omitempty"`
	TabTitleIF *title `json:",omitempty"`
	TabTitleDF *title `json:",omitempty"`

	// PSusers []string `json:omitempty`

	HideconfigMEM *Boole `json:",omitempty"`
	HideconfigIF  *Boole `json:",omitempty"`
	HideconfigCPU *Boole `json:",omitempty"`
	HideconfigDF  *Boole `json:",omitempty"`
	HideconfigPS  *Boole `json:",omitempty"`
	HideconfigVG  *Boole `json:",omitempty"`
}

// server side full client state
type client struct {
	internalClient `json:"-"` // NB not marshalled
	commonClient

	RefreshMEM *refresh `json:",omitempty"`
	RefreshIF  *refresh `json:",omitempty"`
	RefreshCPU *refresh `json:",omitempty"`
	RefreshDF  *refresh `json:",omitempty"`
	RefreshPS  *refresh `json:",omitempty"`
	RefreshVG  *refresh `json:",omitempty"`

	PSplusText       *string `json:",omitempty"`
	PSnotExpandable  *bool   `json:",omitempty"`
	PSnotDecreasable *bool   `json:",omitempty"`
}

type sendClient struct {
	client

	RefreshErrorMEM *bool `json:",omitempty"`
	RefreshErrorIF  *bool `json:",omitempty"`
	RefreshErrorCPU *bool `json:",omitempty"`
	RefreshErrorDF  *bool `json:",omitempty"`
	RefreshErrorPS  *bool `json:",omitempty"`
	RefreshErrorVG  *bool `json:",omitempty"`

	DebugError *string  `json:",omitempty"`
}

type Boole bool
func (b *Boole) merge(src *Boole, send **Boole) {
	if src == nil {
		return
	}
	*b = *src
	*send = src
}


func (_ client) mergeSEQ(dst, src *types.SEQ, send **types.SEQ) {
	if src == nil {
		return
	}
	*dst = *src
	*send = src
}

func(cs *client) Merge(rc recvClient, sc *sendClient) {
	cs.HideMEM.merge(rc.HideMEM, &sc.HideMEM)
	cs.HideIF .merge(rc.HideIF,  &sc.HideIF)
	cs.HideCPU.merge(rc.HideCPU, &sc.HideCPU)
	cs.HideDF .merge(rc.HideDF,  &sc.HideDF)
	cs.HidePS .merge(rc.HidePS,  &sc.HidePS)
	cs.HideVG .merge(rc.HideVG,  &sc.HideVG)

	cs.HideSWAP .merge(rc.HideSWAP,  &sc.HideSWAP)
	cs.ExpandIF .merge(rc.ExpandIF,  &sc.ExpandIF)
	cs.ExpandCPU.merge(rc.ExpandCPU, &sc.ExpandCPU)
	cs.ExpandDF .merge(rc.ExpandDF,  &sc.ExpandDF)

	cs.mergeSEQ(cs.TabIF, rc.TabIF, &sc.TabIF)
	cs.mergeSEQ(cs.TabDF, rc.TabDF, &sc.TabDF)

	cs.TabTitleIF.merge(IFTABS.Title(*cs.TabIF), &sc.TabTitleIF)
	cs.TabTitleDF.merge(DFTABS.Title(*cs.TabDF), &sc.TabTitleDF)

	cs.HideconfigMEM.merge(rc.HideconfigMEM, &sc.HideconfigMEM)
	cs.HideconfigIF .merge(rc.HideconfigIF,  &sc.HideconfigIF)
	cs.HideconfigCPU.merge(rc.HideconfigCPU, &sc.HideconfigCPU)
	cs.HideconfigDF .merge(rc.HideconfigDF,  &sc.HideconfigDF)
	cs.HideconfigPS .merge(rc.HideconfigPS,  &sc.HideconfigPS)
	cs.HideconfigVG .merge(rc.HideconfigVG,  &sc.HideconfigVG)
}

func newtitle(s string) *title {
	p := new(title)
	*p = title(s)
	return p
}

func newfalse()  *bool  { return new(bool); }
func newtrue()   *bool  { return newbool(true); }

func newfalsee() *Boole { return new(Boole); }
func newtruee()  *Boole { return newboole(true); }

func newbool (v bool) (b *bool)  { b = new(bool);  *b = v; return }
func newboole(v bool) (b *Boole) { b = new(Boole); *b = Boole(v); return }

func newseq(v types.SEQ) *types.SEQ {
	s := new(types.SEQ)
	*s = v
	return s
}

func newdefaultrefresh() *refresh {
	r := new(refresh)
	*r = refresh{Duration: periodFlag.Duration}
	return r
}

func defaultClient() client {
	cs := client{}

	cs.HideMEM = newfalsee()
	cs.HideIF  = newfalsee()
	cs.HideCPU = newfalsee()
	cs.HideDF  = newfalsee()
	cs.HidePS  = newfalsee()
	cs.HideVG  = newfalsee()

	cs.HideSWAP  = newfalsee()
	cs.ExpandIF  = newfalsee()
	cs.ExpandCPU = newfalsee()
	cs.ExpandDF  = newfalsee()

	cs.TabIF = newseq(IFBYTES_TABID)
	cs.TabDF = newseq(DFBYTES_TABID)
	cs.TabTitleIF = newtitle(IFTABS.Title(*cs.TabIF))
	cs.TabTitleDF = newtitle(DFTABS.Title(*cs.TabDF))

	hideconfig := true
	// hideconfig  = false // DEVELOPMENT

	cs.HideconfigMEM = newboole(hideconfig)
	cs.HideconfigIF  = newboole(hideconfig)
	cs.HideconfigCPU = newboole(hideconfig)
	cs.HideconfigDF  = newboole(hideconfig)
	cs.HideconfigPS  = newboole(hideconfig)
	cs.HideconfigVG  = newboole(hideconfig)

	cs.RefreshMEM = newdefaultrefresh()
	cs.RefreshIF  = newdefaultrefresh()
	cs.RefreshCPU = newdefaultrefresh()
	cs.RefreshDF  = newdefaultrefresh()
	cs.RefreshPS  = newdefaultrefresh()
	cs.RefreshVG  = newdefaultrefresh()

	cs.psLimit = 8
	// cs.psNotexpandable = newfalse()

	cs.psSEQ = _PSBIMAP.Default_seq
	cs.dfSEQ = _DFBIMAP.Default_seq

	return cs
}
