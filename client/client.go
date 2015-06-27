package client

import (
	"time"

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
		c.RefreshHN,
		c.RefreshUP,
		c.RefreshIP,
		c.RefreshLA,
	}
}

type internalClient struct {
	Toprows int

	MergeRSError error

	Params *Params

	Modified bool
}

type commonClient struct {
}

// server side full client state
type Client struct {
	internalClient `json:"-"` // NB not marshalled
	commonClient

	// un-mergable and hidden refreshes:
	RefreshHN *Refresh `json:"-"`
	RefreshUP *Refresh `json:"-"`
	RefreshIP *Refresh `json:"-"`
	RefreshLA *Refresh `json:"-"`
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

func (c *Client) Merge(r RecvClient, s *SendClient) {
}

// NewClient construct a Client with defaults.
func NewClient(minperiod flags.Period) Client {
	cs := Client{}

	newref := NewRefreshFunc(minperiod)
	// immutable refreshes:
	cs.RefreshHN = newref()
	cs.RefreshUP = newref()
	cs.RefreshIP = newref()
	cs.RefreshLA = newref()

	cs.Params = NewParams(minperiod)
	cs.RecalcRows() // after params

	return cs
}

type RecvClient struct {
	commonClient
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
	// Refresh{HN,UP,IP,LA} are not merged

	err := send.MergeRSError
	send.MergeRSError = nil
	return err
}
