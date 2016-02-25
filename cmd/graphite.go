package cmd

import (
	"log"
	"net"
	"os"
	"time"

	graphite "github.com/cyberdelia/go-metrics-graphite"

	"github.com/ostrost/ostent/ostent"
	"github.com/ostrost/ostent/params"
)

func GraphiteRun(elisting *ostent.ExportingListing, gends params.GraphiteEndpoints) error {
	for _, value := range gends.Values {
		if value.ServerAddr.Host != "" {
			elisting.AddExporter("Graphite", value)
			ostent.AddBackground(GraphiteRunFunc(value))
		}
	}
	return nil
}

func GraphiteRunFunc(value params.Endpoint) func() {
	return func() {
		ostent.Register <- &Carbond{
			ServerAddr: value.ServerAddr.String(),
			Delay:      &value.Delay,
		}
	}
}

type Carbond struct {
	ServerAddr    string
	Conn          net.Conn
	*params.Delay // Expired, Tick methods
	Failing       bool
}

func (cd Carbond) CloseChans()     {} // intentionally empty
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
			grLog.Printf("Sending to %s: %s\n", addr, err)
		}
	} else if cd.Failing {
		cd.Failing = false
		grLog.Printf("%s: Recovered\n", addr)
	}
}

var grLog = log.New(os.Stderr, "[ostent graphite] ", log.LstdFlags)
