package commands

import (
	"flag"
	"net"
	"time"

	graphite "github.com/cyberdelia/go-metrics-graphite"

	"github.com/ostrost/ostent/commands/extpoints"
	"github.com/ostrost/ostent/flags"
	"github.com/ostrost/ostent/ostent"
	"github.com/ostrost/ostent/params"
)

type Graphite struct {
	Log        *extpoints.Log
	DelayFlag  flags.Delay
	ServerAddr flags.Bind
}

func (_ Graphites) SetupFlagSet(cli *flag.FlagSet) extpoints.CommandLineHandler {
	gr := &Graphite{
		Log:        NewLog("[ostent graphite] "),
		DelayFlag:  flags.Delay{Duration: 10 * time.Second}, // 10s default
		ServerAddr: flags.NewBind(2003),
	}
	cli.Var(&gr.DelayFlag, "graphite-delay", "Graphite `delay`")
	cli.Var(&gr.ServerAddr, "graphite-host", "Graphite `host`")
	return func() (extpoints.AtexitHandler, bool, error) {
		if gr.ServerAddr.Host == "" {
			return nil, false, nil
		}
		ostent.AddBackground(func() {
			gc := &Carbond{
				Log:        gr.Log,
				ServerAddr: gr.ServerAddr.String(),
				Delay:      &params.Delay{D: gr.DelayFlag.Duration},
			}
			ostent.Register <- gc
		})
		return nil, false, nil
	}
}

type Carbond struct {
	Log           *extpoints.Log
	ServerAddr    string
	Conn          net.Conn
	*params.Delay // Expired, Tick methods
	Failing       bool
}

func (cd *Carbond) CloseChans()              {} // intentionally empty
func (cd *Carbond) Reload()                  {} // intentionally empty
func (cd *Carbond) Push(*ostent.IndexUpdate) {} // TODO?

func (cd *Carbond) Tack() {
	addr, err := net.ResolveTCPAddr("tcp", cd.ServerAddr)
	if err != nil {
		cd.Log.Printf("Resolve Addr %s: %s\n", cd.ServerAddr, err)
		return
	}
	// go graphite.Graphite(ostent.Reg1s.Registry, 1*time.Second, "ostent", addr)
	err = graphite.GraphiteOnce(graphite.GraphiteConfig{
		DurationUnit: time.Nanosecond, // default, used(divided by thus must not be 0) with Timer metrics
		Addr:         addr,
		Registry:     ostent.Reg1s.Registry,
		Prefix:       "ostent",
	})
	if err != nil {
		if !cd.Failing {
			cd.Failing = true
			cd.Log.Printf("Sending: %s\n", err)
		}
	} else if cd.Failing {
		cd.Failing = false
		cd.Log.Printf("Recovered\n")
	}
}

type Graphites struct{}

func init() {
	extpoints.CommandLines.Register(Graphites{}, "graphite")
}
