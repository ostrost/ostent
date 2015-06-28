package client

import "github.com/ostrost/ostent/flags"

type RecvClient struct{}

// server side full client state
type Client struct {
	Params   *Params `json:"-"`
	Modified bool    `json:"-"`
}

type SendClient struct{ Client }

// NewClient construct a Client with defaults.
func NewClient(minperiod flags.Period) Client {
	return Client{Params: NewParams(minperiod)}
}
