package commands

import (
	"flag"
	"net"
	"time"

	graphite "github.com/cyberdelia/go-metrics-graphite"
	"github.com/ostrost/ostent/flags"
	"github.com/ostrost/ostent/ostent"
	"github.com/ostrost/ostent/params"
)

type Graphite struct {
	Logger      *Logger
	RefreshFlag flags.Period
	ServerAddr  flags.Bind
}

func graphiteCommandLine(cli *flag.FlagSet) CommandLineHandler {
	gr := &Graphite{
		Logger:      NewLogger("[ostent sendto-graphite] "),
		RefreshFlag: flags.Period{Duration: 10 * time.Second}, // 10s default
		ServerAddr:  flags.NewBind(2003),
	}
	cli.Var(&gr.RefreshFlag, "graphite-delay", "Graphite `delay`")
	cli.Var(&gr.ServerAddr, "graphite-host", "Graphite `host`")
	return func() (AtexitHandler, bool, error) {
		if gr.ServerAddr.Host == "" {
			return nil, false, nil
		}
		ostent.AddBackground(func(defaultPeriod flags.Period) {
			/* if gr.RefreshFlag.Duration == 0 { // if .RefreshFlag had no default
				gr.RefreshFlag = defaultPeriod
			} */
			gc := &Carbond{
				Logger:     gr.Logger,
				ServerAddr: gr.ServerAddr.String(),
				Ticks:      params.NewTicks(&params.Duration{D: gr.RefreshFlag.Duration}),
			}
			ostent.Register <- gc
		})
		return nil, false, nil
	}
}

type Carbond struct {
	Logger       *Logger
	ServerAddr   string
	Conn         net.Conn
	params.Ticks // Expired, Tick methods
	Failing      bool
}

func (cd *Carbond) CloseChans()              {} // intentionally empty
func (cd *Carbond) Reload()                  {} // intentionally empty
func (cd *Carbond) Push(*ostent.IndexUpdate) {} // TODO?

func (cd *Carbond) Tack() {
	addr, err := net.ResolveTCPAddr("tcp", cd.ServerAddr)
	if err != nil {
		cd.Logger.Printf("Resolve Addr %s: %s\n", cd.ServerAddr, err)
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
			cd.Logger.Printf("Sending: %s\n", err)
		}
	} else if cd.Failing {
		cd.Failing = false
		cd.Logger.Printf("Recovered\n")
	}
}

func init() {
	AddCommandLine(graphiteCommandLine)
}
