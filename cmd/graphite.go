package cmd

import (
	"net"
	"time"

	graphite "github.com/cyberdelia/go-metrics-graphite"

	"github.com/ostrost/ostent/flags"
	"github.com/ostrost/ostent/ostent"
	"github.com/ostrost/ostent/params"
)

type Graphite struct {
	DelayFlag  flags.Delay
	ServerAddr flags.Bind
}

func (gr *Graphite) Run() error {
	if gr.ServerAddr.Host != "" {
		ostent.AddBackground(func() {
			ostent.AllExporters.AddExporter("graphite")
			ostent.Register <- &Carbond{
				ServerAddr: gr.ServerAddr.String(),
				Delay:      &params.Delay{D: gr.DelayFlag.Duration},
			}
		})
	}
	return nil
}

type Carbond struct {
	ServerAddr    string
	Conn          net.Conn
	*params.Delay // Expired, Tick methods
	Failing       bool
}

func (cd Carbond) CloseChans()     {} // intentionally empty
func (cd Carbond) Reload()         {} // intentionally empty
func (cd Carbond) WantProcs() bool { return false }

func (cd *Carbond) Tack() {
	addr, err := net.ResolveTCPAddr("tcp", cd.ServerAddr)
	if err != nil {
		grLog.Printf("Resolve Addr %s: %s\n", cd.ServerAddr, err)
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
			grLog.Printf("Sending: %s\n", err)
		}
	} else if cd.Failing {
		cd.Failing = false
		grLog.Printf("Recovered\n")
	}
}

var grLog = NewLog(nil, "[ostent graphite] ")
