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
		c.RefreshHN,
		c.RefreshUP,
		c.RefreshIP,
		c.RefreshLA,
	}
}

type internalClient struct {
	PSlimit int

	PSSEQ SEQ
	DFSEQ SEQ

	Toprows int

	Modified bool
}

func (c *Client) mergeTitle(dst *string, src string, send **string) {
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
	c.Modified = true
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

	// un-mergable and hidden refreshes:
	RefreshHN *Refresh `json:"-"`
	RefreshUP *Refresh `json:"-"`
	RefreshIP *Refresh `json:"-"`
	RefreshLA *Refresh `json:"-"`

	PSplusText       *string `json:",omitempty"`
	PSnotExpandable  *bool   `json:",omitempty"`
	PSnotDecreasable *bool   `json:",omitempty"`
}

func (c *Client) RecalcRows() {
	c.Toprows = map[bool]int{true: 1, false: 2}[bool(*c.HideSWAP)]
}

func (sc *SendClient) SetBool(sendb, b **bool, v bool) {
	if Setbool(sendb, b, v) {
		sc.Modified = true
	}
}

func (sc *SendClient) SetString(sends, s **string, v string) {
	if Setstring(sends, s, v) {
		sc.Modified = true
	}
}

func Setbool(sendb, b **bool, v bool) bool {
	if *b != nil && **b == v {
		return false // unchanged
	}
	if *b == nil {
		*b = new(bool)
	}
	**b = v
	*sendb = *b
	return true
}

func Setstring(sends, s **string, v string) bool {
	if *s != nil && **s == v {
		return false // unchanged
	}
	if *s == nil {
		*s = new(string)
	}
	**s = v
	*sends = *s
	return true
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

func (c *Client) mergeBool(dst, src *bool, send **bool) {
	// c is unused
	if src == nil {
		return
	}
	*dst = *src
	*send = src
	c.Modified = true
}

func (c *Client) mergeSEQ(dst, src *SEQ, send **SEQ) {
	// c is unused
	if src == nil {
		return
	}
	*dst = *src
	*send = src
	c.Modified = true
}

func (c *Client) Merge(r RecvClient, s *SendClient) {
	s.mergeBool(c.HideRAM, r.HideRAM, &s.HideRAM)
	s.mergeBool(c.HideIF, r.HideIF, &s.HideIF)
	s.mergeBool(c.HideCPU, r.HideCPU, &s.HideCPU)
	s.mergeBool(c.HideDF, r.HideDF, &s.HideDF)
	s.mergeBool(c.HidePS, r.HidePS, &s.HidePS)
	s.mergeBool(c.HideVG, r.HideVG, &s.HideVG)

	s.mergeBool(c.HideSWAP, r.HideSWAP, &s.HideSWAP)
	s.mergeBool(c.ExpandIF, r.ExpandIF, &s.ExpandIF)
	s.mergeBool(c.ExpandCPU, r.ExpandCPU, &s.ExpandCPU)
	s.mergeBool(c.ExpandDF, r.ExpandDF, &s.ExpandDF)

	s.mergeBool(c.HideconfigMEM, r.HideconfigMEM, &s.HideconfigMEM)
	s.mergeBool(c.HideconfigIF, r.HideconfigIF, &s.HideconfigIF)
	s.mergeBool(c.HideconfigCPU, r.HideconfigCPU, &s.HideconfigCPU)
	s.mergeBool(c.HideconfigDF, r.HideconfigDF, &s.HideconfigDF)
	s.mergeBool(c.HideconfigPS, r.HideconfigPS, &s.HideconfigPS)
	s.mergeBool(c.HideconfigVG, r.HideconfigVG, &s.HideconfigVG)

	s.mergeSEQ(c.TabIF, r.TabIF, &s.TabIF)
	s.mergeSEQ(c.TabDF, r.TabDF, &s.TabDF)

	// merge NOT from the r
	s.mergeTitle(c.TabTitleIF, IFTABS.Title(*c.TabIF), &s.TabTitleIF)
	s.mergeTitle(c.TabTitleDF, DFTABS.Title(*c.TabDF), &s.TabTitleDF)
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

	// immutable refreshes:
	cs.RefreshHN = &Refresh{Period: minperiod}
	cs.RefreshUP = &Refresh{Period: minperiod}
	cs.RefreshIP = &Refresh{Period: minperiod}
	cs.RefreshLA = &Refresh{Period: minperiod}

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
	// cs modification of .PSlimit does not flip .Modified
	if *rs.MorePsignal {
		if cs.PSlimit < 65536 {
			cs.PSlimit *= 2
		}
	} else if cs.PSlimit >= 2 {
		cs.PSlimit /= 2
	}
	rs.MorePsignal = nil
}

func (sc *SendClient) mergeRefreshSignal(above time.Duration, ppinput *string, prefresh *Refresh, sendr **Refresh, senderr **bool) error {
	if ppinput == nil {
		return nil
	}
	pv := flags.Period{Above: &above}
	if err := pv.Set(*ppinput); err != nil {
		*senderr = newtrue()
		sc.Modified = true // otherwise refresh input error won't be sent
		return err
	}
	*senderr = newfalse()
	*sendr = new(Refresh)
	(**sendr).Duration = pv.Duration
	prefresh.Duration = pv.Duration
	prefresh.tick = 0
	sc.Modified = true
	return nil
}

func (rs *RecvClient) MergeClient(minperiod flags.Period, cs *Client, send *SendClient) error {
	minrefresh := minperiod.Duration
	rs.mergeMorePsignal(cs)
	var refreshmem Refresh
	if err := send.mergeRefreshSignal(minrefresh, rs.RefreshSignalMEM, &refreshmem, &send.RefreshMEM, &send.RefreshErrorMEM); err != nil {
		return err
	}
	// refreshmem value change
	*cs.RefreshRAM = refreshmem
	*cs.RefreshSWAP = refreshmem

	if err := send.mergeRefreshSignal(minrefresh, rs.RefreshSignalIF, cs.RefreshIF, &send.RefreshIF, &send.RefreshErrorIF); err != nil {
		return err
	}
	if err := send.mergeRefreshSignal(minrefresh, rs.RefreshSignalCPU, cs.RefreshCPU, &send.RefreshCPU, &send.RefreshErrorCPU); err != nil {
		return err
	}
	if err := send.mergeRefreshSignal(minrefresh, rs.RefreshSignalDF, cs.RefreshDF, &send.RefreshDF, &send.RefreshErrorDF); err != nil {
		return err
	}
	if err := send.mergeRefreshSignal(minrefresh, rs.RefreshSignalPS, cs.RefreshPS, &send.RefreshPS, &send.RefreshErrorPS); err != nil {
		return err
	}
	if err := send.mergeRefreshSignal(minrefresh, rs.RefreshSignalVG, cs.RefreshVG, &send.RefreshVG, &send.RefreshErrorVG); err != nil {
		return err
	}
	// Refresh{HN,UP,IP,LA} are not merged
	return nil
}
