package commands

import (
	"flag"
	"time"

	librato "github.com/mihasya/go-metrics-librato"

	"github.com/ostrost/ostent/commands/extpoints"
	"github.com/ostrost/ostent/flags"
	"github.com/ostrost/ostent/ostent"
)

type Librato struct {
	DelayFlag            flags.Delay
	Email, Token, Source string
}

func (_ Libratos) SetupFlagSet(cli *flag.FlagSet) extpoints.CommandLineHandler {
	hostname, err := ostent.GetHN()
	if err != nil {
		hostname = ""
	}
	lr := &Librato{
		DelayFlag: flags.Delay{Duration: 10 * time.Second}, // 10s default
	}
	cli.Var(&lr.DelayFlag, "librato-delay", "Librato `delay`")
	cli.StringVar(&lr.Email, "librato-email", "", "Librato `email`")
	cli.StringVar(&lr.Token, "librato-token", "", "Librato `token`")
	cli.StringVar(&lr.Source, "librato-source", hostname, "Librato `source`")
	return func() (extpoints.AtexitHandler, bool, error) {
		if lr.Email == "" {
			return nil, false, nil
		}
		ostent.AddBackground(func() {
			go librato.Librato(ostent.Reg1s.Registry,
				lr.DelayFlag.Duration,
				lr.Email,
				lr.Token,
				lr.Source,
				[]float64{0.95},
				time.Millisecond,
			)
		})
		return nil, false, nil
	}
}

type Libratos struct{}

func init() {
	extpoints.CommandLines.Register(Libratos{}, "librato")
}
