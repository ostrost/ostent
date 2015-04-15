// Package client is all about client state.
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

	PSSEQ Number
	DFSEQ Number

	Toprows int

	MergeRSError error

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

	TabIF *Uint `json:",omitempty"`
	TabDF *Uint `json:",omitempty"`

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

func (c *Client) mergeUint(dst, src *Uint, send **Uint) {
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

	s.mergeUint(c.TabIF, r.TabIF, &s.TabIF)
	s.mergeUint(c.TabDF, r.TabDF, &s.TabDF)

	// merge NOT from the r
	s.mergeTitle(c.TabTitleIF, IFTABS.Title(*c.TabIF), &s.TabTitleIF)
	s.mergeTitle(c.TabTitleDF, DFTABS.Title(*c.TabDF), &s.TabTitleDF)
}

func DefaultClient(minperiod flags.Period) Client {
	cs := Client{}

	// new(bool) is &false
	cs.HideRAM = new(bool)
	cs.HideIF = new(bool)
	cs.HideCPU = new(bool)
	cs.HideDF = new(bool)
	cs.HidePS = new(bool)
	cs.HideVG = new(bool)
	cs.HideSWAP = new(bool)
	cs.ExpandIF = new(bool)
	cs.ExpandCPU = new(bool)
	cs.ExpandDF = new(bool)

	newuint := func(p Uint) *Uint { n := new(Uint); *n = p; return n }
	cs.TabIF = newuint(IFBYTES_TABID)
	cs.TabDF = newuint(DFBYTES_TABID)

	cs.TabTitleIF = new(string)
	*cs.TabTitleIF = IFTABS.Title(*cs.TabIF)
	cs.TabTitleDF = new(string)
	*cs.TabTitleDF = DFTABS.Title(*cs.TabDF)

	newhc := func() *bool { b := new(bool); *b = true; return b } // *b = false for DEVELOPMENT
	cs.HideconfigMEM = newhc()
	cs.HideconfigIF = newhc()
	cs.HideconfigCPU = newhc()
	cs.HideconfigDF = newhc()
	cs.HideconfigPS = newhc()
	cs.HideconfigVG = newhc()

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

	cs.PSSEQ = PS.Default
	cs.DFSEQ = DF.Default

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

// mergeRefreshSignal returns true when prefresh is modified.
func (sc *SendClient) mergeRefreshSignal(above time.Duration, ppinput *string, prefresh *Refresh, sendr **Refresh, senderr **bool) bool {
	if sc.MergeRSError != nil {
		return false
	}
	if ppinput == nil {
		return false
	}
	*senderr = new(bool) // false by default
	sc.Modified = true   // senderr is non-nil, ergo sc is modified
	pv := flags.Period{Above: &above}
	if err := pv.Set(*ppinput); err != nil {
		sc.MergeRSError = err
		**senderr = true
		return false
	}
	*sendr = new(Refresh)
	(**sendr).Duration = pv.Duration
	prefresh.Duration = pv.Duration
	prefresh.tick = 0
	return true
}

// MergeRefresh merges into cs various refresh updates. send is populated with the updates.
func (rs *RecvClient) MergeRefresh(minrefresh time.Duration, cs *Client, send *SendClient) error {
	rs.mergeMorePsignal(cs)

	rrammod := send.mergeRefreshSignal(minrefresh, rs.RefreshSignalMEM, cs.RefreshRAM, &send.RefreshMEM, &send.RefreshErrorMEM)
	if send.MergeRSError == nil && rrammod { // RefreshRAM value change, so should RefreshSWAP
		*cs.RefreshSWAP = *cs.RefreshRAM
	}

	send.mergeRefreshSignal(minrefresh, rs.RefreshSignalIF, cs.RefreshIF, &send.RefreshIF, &send.RefreshErrorIF)
	send.mergeRefreshSignal(minrefresh, rs.RefreshSignalCPU, cs.RefreshCPU, &send.RefreshCPU, &send.RefreshErrorCPU)
	send.mergeRefreshSignal(minrefresh, rs.RefreshSignalDF, cs.RefreshDF, &send.RefreshDF, &send.RefreshErrorDF)
	send.mergeRefreshSignal(minrefresh, rs.RefreshSignalPS, cs.RefreshPS, &send.RefreshPS, &send.RefreshErrorPS)
	send.mergeRefreshSignal(minrefresh, rs.RefreshSignalVG, cs.RefreshVG, &send.RefreshVG, &send.RefreshErrorVG)
	// Refresh{HN,UP,IP,LA} are not merged

	err := send.MergeRSError
	send.MergeRSError = nil
	return err
}
