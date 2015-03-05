package client

import (
	"time"

	"github.com/ostrost/ostent/flags"
)

type Refresh struct {
	flags.Period
	tick int // .Tick() must be called once per second; .tick is 1 when the refresh expired
}

func (r *Refresh) Refresh(forcerefresh bool) bool {
	if forcerefresh {
		return true
	}
	return r.expired()
}

func (r Refresh) expired() bool {
	return r.tick <= 1
}

func (c *Client) Tick() {
	for _, r := range c.refreshes() {
		r.tick++
		if r.tick-1 >= int(time.Duration(r.Duration)/time.Second) {
			r.tick = 1 // expired
		}
	}
}

func (c Client) Expired() bool {
	for _, r := range c.refreshes() {
		if r.expired() {
			return true
		}
	}
	return false
}

func (c *Client) refreshes() []*Refresh {
	return []*Refresh{
		c.RefreshRAM,
		c.RefreshSWAP,
		c.RefreshIF,
		c.RefreshCPU,
		c.RefreshDF,
		c.RefreshPS,
		c.RefreshVG,
	}
}

type internalClient struct {
	// NB lowercase fields only, NOT to be marshalled/exported

	PSlimit int

	PSSEQ SEQ
	DFSEQ SEQ

	Toprows int
}

func (c Client) mergeTitle(dst *string, src string, send **string) {
	if src == "" { // precautious. should not be the case
		return
	}
	// *send = nil
	if *dst == src {
		return
	}
	*send = new(string)
	**send = src
	*dst = **send
}

type commonClient struct {
	HideRAM   *bool `json:",omitempty"`
	HideIF    *bool `json:",omitempty"`
	HideCPU   *bool `json:",omitempty"`
	HideDF    *bool `json:",omitempty"`
	HidePS    *bool `json:",omitempty"`
	HideVG    *bool `json:",omitempty"`
	HideSWAP  *bool `json:",omitempty"`
	ExpandIF  *bool `json:",omitempty"`
	ExpandCPU *bool `json:",omitempty"`
	ExpandDF  *bool `json:",omitempty"`

	TabIF *SEQ `json:",omitempty"`
	TabDF *SEQ `json:",omitempty"`

	TabTitleIF *string `json:",omitempty"`
	TabTitleDF *string `json:",omitempty"`

	// PSusers []string `json:omitempty`

	HideconfigMEM *bool `json:",omitempty"`
	HideconfigIF  *bool `json:",omitempty"`
	HideconfigCPU *bool `json:",omitempty"`
	HideconfigDF  *bool `json:",omitempty"`
	HideconfigPS  *bool `json:",omitempty"`
	HideconfigVG  *bool `json:",omitempty"`
}

// server side full client state
type Client struct {
	internalClient `json:"-"` // NB not marshalled
	commonClient

	ExpandableIF  *bool `json:",omitempty"`
	ExpandableCPU *bool `json:",omitempty"`
	ExpandableDF  *bool `json:",omitempty"`

	ExpandtextIF  *string `json:",omitempty"`
	ExpandtextCPU *string `json:",omitempty"`
	ExpandtextDF  *string `json:",omitempty"`

	// RefreshGeneric *refresh `json:",omitempty"`
	RefreshRAM  *Refresh `json:",omitempty"`
	RefreshSWAP *Refresh `json:",omitempty"`
	RefreshIF   *Refresh `json:",omitempty"`
	RefreshCPU  *Refresh `json:",omitempty"`
	RefreshDF   *Refresh `json:",omitempty"`
	RefreshPS   *Refresh `json:",omitempty"`
	RefreshVG   *Refresh `json:",omitempty"`

	PSplusText       *string `json:",omitempty"`
	PSnotExpandable  *bool   `json:",omitempty"`
	PSnotDecreasable *bool   `json:",omitempty"`
}

func (c *Client) RecalcRows() {
	c.Toprows = map[bool]int{true: 1, false: 2}[bool(*c.HideSWAP)]
}

func SetBool(b, b2 **bool, v bool) {
	if *b != nil && **b == v {
		return // unchanged
	}
	if *b == nil {
		*b = new(bool)
	}
	**b = v
	*b2 = *b
}

func SetString(s, s2 **string, v string) {
	if *s != nil && **s == v {
		return // unchanged
	}
	if *s == nil {
		*s = new(string)
	}
	**s = v
	*s2 = *s
}

type SendClient struct {
	Client

	RefreshErrorMEM  *bool `json:",omitempty"`
	RefreshErrorSWAP *bool `json:",omitempty"`
	RefreshErrorIF   *bool `json:",omitempty"`
	RefreshErrorCPU  *bool `json:",omitempty"`
	RefreshErrorDF   *bool `json:",omitempty"`
	RefreshErrorPS   *bool `json:",omitempty"`
	RefreshErrorVG   *bool `json:",omitempty"`

	RefreshMEM  *Refresh  `json:",omitempty"`  // for frontend only
	RefreshRAM  *struct{} `json:"-,omitempty"` // shadow
	RefreshSWAP *struct{} `json:"-,omitempty"` // shadow

	DebugError *string `json:",omitempty"`
}

func (c Client) mergeBool(dst, src *bool, send **bool) {
	// c is unused
	if src == nil {
		return
	}
	*dst = *src
	*send = src
}

func (c Client) mergeSEQ(dst, src *SEQ, send **SEQ) {
	// c is unused
	if src == nil {
		return
	}
	*dst = *src
	*send = src
}

func (c *Client) Merge(r RecvClient, s *SendClient) {
	c.mergeBool(c.HideRAM, r.HideRAM, &s.HideRAM)
	c.mergeBool(c.HideIF, r.HideIF, &s.HideIF)
	c.mergeBool(c.HideCPU, r.HideCPU, &s.HideCPU)
	c.mergeBool(c.HideDF, r.HideDF, &s.HideDF)
	c.mergeBool(c.HidePS, r.HidePS, &s.HidePS)
	c.mergeBool(c.HideVG, r.HideVG, &s.HideVG)

	c.mergeBool(c.HideSWAP, r.HideSWAP, &s.HideSWAP)
	c.mergeBool(c.ExpandIF, r.ExpandIF, &s.ExpandIF)
	c.mergeBool(c.ExpandCPU, r.ExpandCPU, &s.ExpandCPU)
	c.mergeBool(c.ExpandDF, r.ExpandDF, &s.ExpandDF)

	c.mergeBool(c.HideconfigMEM, r.HideconfigMEM, &s.HideconfigMEM)
	c.mergeBool(c.HideconfigIF, r.HideconfigIF, &s.HideconfigIF)
	c.mergeBool(c.HideconfigCPU, r.HideconfigCPU, &s.HideconfigCPU)
	c.mergeBool(c.HideconfigDF, r.HideconfigDF, &s.HideconfigDF)
	c.mergeBool(c.HideconfigPS, r.HideconfigPS, &s.HideconfigPS)
	c.mergeBool(c.HideconfigVG, r.HideconfigVG, &s.HideconfigVG)

	c.mergeSEQ(c.TabIF, r.TabIF, &s.TabIF)
	c.mergeSEQ(c.TabDF, r.TabDF, &s.TabDF)

	// merge NOT from the r
	c.mergeTitle(c.TabTitleIF, IFTABS.Title(*c.TabIF), &s.TabTitleIF)
	c.mergeTitle(c.TabTitleDF, DFTABS.Title(*c.TabDF), &s.TabTitleDF)
}

func newfalse() *bool      { return new(bool) }
func newtrue() *bool       { return newbool(true) }
func newbool(v bool) *bool { b := new(bool); *b = v; return b }

func newseq(v SEQ) *SEQ {
	s := new(SEQ)
	*s = v
	return s
}

func DefaultClient(minperiod flags.Period) Client {
	cs := Client{}

	cs.HideRAM = newfalse()
	cs.HideIF = newfalse()
	cs.HideCPU = newfalse()
	cs.HideDF = newfalse()
	cs.HidePS = newfalse()
	cs.HideVG = newfalse()

	cs.HideSWAP = newfalse()
	cs.ExpandIF = newfalse()
	cs.ExpandCPU = newfalse()
	cs.ExpandDF = newfalse()

	cs.TabIF = newseq(IFBYTES_TABID)
	cs.TabDF = newseq(DFBYTES_TABID)

	cs.TabTitleIF = new(string)
	*cs.TabTitleIF = IFTABS.Title(*cs.TabIF)
	cs.TabTitleDF = new(string)
	*cs.TabTitleDF = DFTABS.Title(*cs.TabDF)

	hideconfig := true
	// hideconfig  = false // DEVELOPMENT

	cs.HideconfigMEM = newbool(hideconfig)
	cs.HideconfigIF = newbool(hideconfig)
	cs.HideconfigCPU = newbool(hideconfig)
	cs.HideconfigDF = newbool(hideconfig)
	cs.HideconfigPS = newbool(hideconfig)
	cs.HideconfigVG = newbool(hideconfig)

	//cs.RefreshGeneric = &refresh{Period: minperiod}
	cs.RefreshRAM = &Refresh{Period: minperiod}
	cs.RefreshSWAP = &Refresh{Period: minperiod}
	cs.RefreshIF = &Refresh{Period: minperiod}
	cs.RefreshCPU = &Refresh{Period: minperiod}
	cs.RefreshDF = &Refresh{Period: minperiod}
	cs.RefreshPS = &Refresh{Period: minperiod}
	cs.RefreshVG = &Refresh{Period: minperiod}

	cs.PSlimit = 8

	cs.PSSEQ = PSBIMAP.DefaultSeq
	cs.DFSEQ = DFBIMAP.DefaultSeq

	cs.RecalcRows()

	return cs
}

type RecvClient struct {
	commonClient
	MorePsignal      *bool
	RefreshSignalMEM *string
	RefreshSignalIF  *string
	RefreshSignalCPU *string
	RefreshSignalDF  *string
	RefreshSignalPS  *string
	RefreshSignalVG  *string
}

func (rs *RecvClient) mergeMorePsignal(cs *Client) {
	if rs.MorePsignal == nil {
		return
	}
	if *rs.MorePsignal {
		if cs.PSlimit < 65536 {
			cs.PSlimit *= 2
		}
	} else if cs.PSlimit >= 2 {
		cs.PSlimit /= 2
	}
	rs.MorePsignal = nil
}

func (rs *RecvClient) mergeRefreshSignal(above time.Duration, ppinput *string, prefresh *Refresh, sendr **Refresh, senderr **bool) error {
	if ppinput == nil {
		return nil
	}
	pv := flags.Period{Above: &above}
	if err := pv.Set(*ppinput); err != nil {
		*senderr = newtrue()
		return err
	}
	*senderr = newfalse()
	*sendr = new(Refresh)
	(**sendr).Duration = pv.Duration
	prefresh.Duration = pv.Duration
	prefresh.tick = 0
	return nil
}

func (rs *RecvClient) MergeClient(minperiod flags.Period, cs *Client, send *SendClient) error {
	minrefresh := minperiod.Duration
	rs.mergeMorePsignal(cs)
	var refreshmem Refresh
	if err := rs.mergeRefreshSignal(minrefresh, rs.RefreshSignalMEM, &refreshmem, &send.RefreshMEM, &send.RefreshErrorMEM); err != nil {
		return err
	} else {
		*cs.RefreshRAM = refreshmem
		*cs.RefreshSWAP = refreshmem
	}
	if err := rs.mergeRefreshSignal(minrefresh, rs.RefreshSignalIF, cs.RefreshIF, &send.RefreshIF, &send.RefreshErrorIF); err != nil {
		return err
	}
	if err := rs.mergeRefreshSignal(minrefresh, rs.RefreshSignalCPU, cs.RefreshCPU, &send.RefreshCPU, &send.RefreshErrorCPU); err != nil {
		return err
	}
	if err := rs.mergeRefreshSignal(minrefresh, rs.RefreshSignalDF, cs.RefreshDF, &send.RefreshDF, &send.RefreshErrorDF); err != nil {
		return err
	}
	if err := rs.mergeRefreshSignal(minrefresh, rs.RefreshSignalPS, cs.RefreshPS, &send.RefreshPS, &send.RefreshErrorPS); err != nil {
		return err
	}
	if err := rs.mergeRefreshSignal(minrefresh, rs.RefreshSignalVG, cs.RefreshVG, &send.RefreshVG, &send.RefreshErrorVG); err != nil {
		return err
	}
	return nil
}
