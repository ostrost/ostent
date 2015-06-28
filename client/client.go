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
	Default flags.Period `json:"-"` // read-only
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

type RecvClient struct{}

type internalClient struct {
	Toprows  int
	Params   *Params
	Modified bool
}

// server side full client state
type Client struct {
	internalClient `json:"-"` // NB not marshalled

	// un-mergable and hidden refreshes:
	RefreshHN *Refresh `json:"-"`
	RefreshUP *Refresh `json:"-"`
	RefreshIP *Refresh `json:"-"`
	RefreshLA *Refresh `json:"-"`
}

func (c *Client) RecalcRows() {
	c.Toprows = map[bool]int{true: 1, false: 2}[c.Params.BOOL["hideswap"].Value]
}

type SendClient struct {
	Client
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
