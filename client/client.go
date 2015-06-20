package client

import (
	"time"

	"github.com/ostrost/ostent/client/enums"
	"github.com/ostrost/ostent/flags"
)

// Refresh is a ticker with period.
// json.Marshal exposes inline .Period only,
// .Default is explicitly ignored.
// .Default is available and used in templates.
type Refresh struct {
	flags.Period
	Default flags.Period `json:"-"` // not modified ever, including in .Merge* funcs
	tick    int          // .Tick() must be called once per second; .tick is 1 when the refresh expired
}

// NewRefreshFunc constructs Refresh-maker.
func NewRefreshFunc(period flags.Period) func() *Refresh {
	return func() *Refresh {
		return &Refresh{Period: period, Default: period}
	}
}

// TODO .Refresh method is used in ostent.Set/Refresher only. To be removed.
func (r *Refresh) Refresh(forcerefresh bool) bool {
	if forcerefresh {
		return true
	}
	return r.tick <= 1 // r.expired()
}

// func (r Refresh) expired() bool { return r.tick <= 1 }

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
		if r.tick <= 1 { // if r.expired()
			return true
		}
	}
	return false
}

func (c *Client) refreshes() []*Refresh {
	return []*Refresh{
		c.RefreshMME, // Used to be RefreshMEM, soon to be gone.
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

	Toprows int

	MergeRSError error

	Params *Params

	Modified bool
}

type commonClient struct {
	HideIF    *bool `json:",omitempty"`
	HideCPU   *bool `json:",omitempty"`
	HideDF    *bool `json:",omitempty"`
	HidePS    *bool `json:",omitempty"`
	HideVG    *bool `json:",omitempty"`
	ExpandIF  *bool `json:",omitempty"`
	ExpandCPU *bool `json:",omitempty"`
	ExpandDF  *bool `json:",omitempty"`

	TabIF *Tab `json:",omitempty"`
	TabDF *Tab `json:",omitempty"`

	// PSusers []string `json:omitempty`

	// HideconfigMEM *bool `json:",omitempty"`
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

	RefreshMME *Refresh `json:",omitempty"` // Used to be RefreshMEM, soon to be gone.
	RefreshIF  *Refresh `json:",omitempty"`
	RefreshCPU *Refresh `json:",omitempty"`
	RefreshDF  *Refresh `json:",omitempty"`
	RefreshPS  *Refresh `json:",omitempty"`
	RefreshVG  *Refresh `json:",omitempty"`

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
	c.Toprows = map[bool]int{true: 1, false: 2}[c.Params.BOOL["hideswap"].Value]
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

	// RefreshErrorMEM  *bool `json:",omitempty"`
	RefreshErrorSWAP *bool `json:",omitempty"`
	RefreshErrorIF   *bool `json:",omitempty"`
	RefreshErrorCPU  *bool `json:",omitempty"`
	RefreshErrorDF   *bool `json:",omitempty"`
	RefreshErrorPS   *bool `json:",omitempty"`
	RefreshErrorVG   *bool `json:",omitempty"`

	DebugError *string `json:",omitempty"`
}

func (c *Client) mergeBool(dst, src *bool, send **bool) {
	if src == nil {
		return
	}
	*dst = *src
	*send = src
	c.Modified = true
}

// MergeTab fills dst and send with src.Uint and title found from tabs.
// src.Title is disregarded.
func (c *Client) MergeTab(dst, src *Tab, send **Tab, tabs Tabs) {
	if src == nil || src.Uint == dst.Uint {
		return
	}
	var title string
	for _, v := range tabs {
		if src.Uint == v.Uint {
			title = v.Title
			break
		}
	}
	if title == "" { // no tab with src.Uint found
		return
	}
	if send == nil {
		s := new(Tab)
		send = &s // dummy
	} else {
		*send = new(Tab)
	}
	dst.Uint, (*send).Uint = src.Uint, src.Uint
	dst.Title, (*send).Title = title, title
	c.Modified = true
}

func (c *Client) NewTab(tabs Tabs, u enums.Uint) *Tab {
	n := new(Tab)
	c.MergeTab(n, &Tab{Uint: u}, nil, tabs)
	return n
}

func (c *Client) Merge(r RecvClient, s *SendClient) {
	s.mergeBool(c.HideIF, r.HideIF, &s.HideIF)
	s.mergeBool(c.HideCPU, r.HideCPU, &s.HideCPU)
	s.mergeBool(c.HideDF, r.HideDF, &s.HideDF)
	s.mergeBool(c.HidePS, r.HidePS, &s.HidePS)
	s.mergeBool(c.HideVG, r.HideVG, &s.HideVG)

	s.mergeBool(c.ExpandIF, r.ExpandIF, &s.ExpandIF)
	s.mergeBool(c.ExpandCPU, r.ExpandCPU, &s.ExpandCPU)
	s.mergeBool(c.ExpandDF, r.ExpandDF, &s.ExpandDF)

	// s.mergeBool(c.HideconfigMEM, r.HideconfigMEM, &s.HideconfigMEM)
	s.mergeBool(c.HideconfigIF, r.HideconfigIF, &s.HideconfigIF)
	s.mergeBool(c.HideconfigCPU, r.HideconfigCPU, &s.HideconfigCPU)
	s.mergeBool(c.HideconfigDF, r.HideconfigDF, &s.HideconfigDF)
	s.mergeBool(c.HideconfigPS, r.HideconfigPS, &s.HideconfigPS)
	s.mergeBool(c.HideconfigVG, r.HideconfigVG, &s.HideconfigVG)

	s.MergeTab(c.TabIF, r.TabIF, &s.TabIF, IFTABS)
	s.MergeTab(c.TabDF, r.TabDF, &s.TabDF, DFTABS)
}

// NewClient construct a Client with defaults.
func NewClient(minperiod flags.Period) Client {
	cs := Client{}

	// new(bool) is &false
	cs.HideIF = new(bool)
	cs.HideCPU = new(bool)
	cs.HideDF = new(bool)
	cs.HidePS = new(bool)
	cs.HideVG = new(bool)
	cs.ExpandIF = new(bool)
	cs.ExpandCPU = new(bool)
	cs.ExpandDF = new(bool)

	cs.TabIF = cs.NewTab(IFTABS, IFBYTES)
	cs.TabDF = cs.NewTab(DFTABS, DFBYTES)

	newhc := func() *bool { b := new(bool); *b = true; return b } // *b = false for DEVELOPMENT
	// cs.HideconfigMEM = newhc()
	cs.HideconfigIF = newhc()
	cs.HideconfigCPU = newhc()
	cs.HideconfigDF = newhc()
	cs.HideconfigPS = newhc()
	cs.HideconfigVG = newhc()

	newref := NewRefreshFunc(minperiod)
	cs.RefreshMME = newref() // Used to be RefreshMEM, soon to be gone.
	cs.RefreshIF = newref()
	cs.RefreshCPU = newref()
	cs.RefreshDF = newref()
	cs.RefreshPS = newref()
	cs.RefreshVG = newref()

	// immutable refreshes:
	cs.RefreshHN = newref()
	cs.RefreshUP = newref()
	cs.RefreshIP = newref()
	cs.RefreshLA = newref()

	cs.PSlimit = 8

	cs.Params = NewParams(minperiod)
	cs.RecalcRows() // after params

	return cs
}

type RecvClient struct {
	commonClient
	MorePsignal      *bool
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

// MergeRefreshSignal stores parsed ppinput into prefresh AND sendr or error in senderr.
func (sc *SendClient) MergeRefreshSignal(ppinput *string, prefresh *Refresh, sendr **Refresh, senderr **bool) {
	if sc.MergeRSError != nil {
		return
	}
	if ppinput == nil {
		return
	}
	*senderr = new(bool) // false by default
	sc.Modified = true   // senderr is non-nil, ergo sc is modified
	pv := flags.Period{Above: &prefresh.Default.Duration}
	if err := pv.Set(*ppinput); err != nil {
		sc.MergeRSError = err
		**senderr = true
		return
	}
	*sendr = new(Refresh)
	(**sendr).Duration = pv.Duration
	prefresh.Duration = pv.Duration
	prefresh.tick = 0
	return
}

// MergeRefresh merges into cs various refresh updates. send is populated with the updates.
func (rs *RecvClient) MergeRefresh(cs *Client, send *SendClient) error {
	rs.mergeMorePsignal(cs)

	send.MergeRefreshSignal(rs.RefreshSignalIF, cs.RefreshIF, &send.RefreshIF, &send.RefreshErrorIF)
	send.MergeRefreshSignal(rs.RefreshSignalCPU, cs.RefreshCPU, &send.RefreshCPU, &send.RefreshErrorCPU)
	send.MergeRefreshSignal(rs.RefreshSignalDF, cs.RefreshDF, &send.RefreshDF, &send.RefreshErrorDF)
	send.MergeRefreshSignal(rs.RefreshSignalPS, cs.RefreshPS, &send.RefreshPS, &send.RefreshErrorPS)
	send.MergeRefreshSignal(rs.RefreshSignalVG, cs.RefreshVG, &send.RefreshVG, &send.RefreshErrorVG)
	// Refresh{HN,UP,IP,LA} are not merged

	err := send.MergeRSError
	send.MergeRSError = nil
	return err
}
