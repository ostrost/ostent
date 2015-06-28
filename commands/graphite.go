package commands

import (
	"flag"
	"net"
	"time"

	"github.com/ostrost/ostent/flags"
	"github.com/ostrost/ostent/ostent"
	"github.com/ostrost/ostent/params"
	metrics "github.com/rcrowley/go-metrics"
)

type graphite struct {
	logger      *Logger
	RefreshFlag flags.Period
	ServerAddr  flags.Bind
}

func graphiteCommandLine(cli *flag.FlagSet) CommandLineHandler {
	gr := &graphite{
		logger:      NewLogger("[ostent sendto-graphite] "),
		RefreshFlag: flags.Period{Duration: 10 * time.Second}, // 10s default
		ServerAddr:  flags.NewBind(2003),
	}
	cli.Var(&gr.RefreshFlag, "graphite-refresh", "Graphite refresh interval")
	cli.Var(&gr.ServerAddr, "sendto-graphite", "Graphite server address")
	return func() (AtexitHandler, bool, error) {
		if gr.ServerAddr.Host == "" {
			return nil, false, nil
		}
		ostent.AddBackground(func(defaultPeriod flags.Period) {
			/* if gr.RefreshFlag.Duration == 0 { // if .RefreshFlag had no default
				gr.RefreshFlag = defaultPeriod
			} */
			gc := &carbond{
				logger:      gr.logger,
				serveraddr:  gr.ServerAddr.String(),
				PeriodParam: params.NewPeriodParam(gr.RefreshFlag, "refreshgraphite", nil),
			}
			ostent.Register <- gc
		})
		return nil, false, nil
	}
}

type carbond struct {
	logger     *Logger
	serveraddr string
	conn       net.Conn
	*params.PeriodParam
	failing bool
}

func (_ *carbond) CloseChans() {} // intentionally empty
func (_ *carbond) Reload()     {} // intentionally empty

func (_ *carbond) Push(*ostent.IndexUpdate) {} // TODO?

func (cd *carbond) Tack() {
	addr, err := net.ResolveTCPAddr("tcp", cd.serveraddr)
	if err != nil {
		cd.logger.Printf("Resolve Addr %s: %s\n", cd.serveraddr, err)
		return
	}
	// go metrics.Graphite(ostent.Reg1s.Registry, 1*time.Second, "ostent", addr)
	err = metrics.GraphiteOnce(metrics.GraphiteConfig{
		DurationUnit: time.Nanosecond, // default, used(divided by thus must not be 0) with Timer metrics
		Addr:         addr,
		Registry:     ostent.Reg1s.Registry,
		Prefix:       "ostent",
	})
	if err != nil {
		if !cd.failing {
			cd.failing = true
			cd.logger.Printf("Sending: %s\n", err)
		}
	} else if cd.failing {
		cd.failing = false
		cd.logger.Printf("Recovered\n")
	}
}

func init() {
	AddCommandLine(graphiteCommandLine)
}
