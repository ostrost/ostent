package cmd

import (
	"time"

	librato "github.com/mihasya/go-metrics-librato"

	"github.com/ostrost/ostent/flags"
	"github.com/ostrost/ostent/ostent"
)

type Librato struct {
	DelayFlag            flags.Delay
	Email, Token, Source string
}

func (lr *Librato) Run() error {
	if lr.Email != "" {
		ostent.AddBackground(func() {
			ostent.AllExporters.AddExporter("librato")
			go librato.Librato(ostent.Reg1s.Registry,
				lr.DelayFlag.Duration,
				lr.Email,
				lr.Token,
				lr.Source,
				[]float64{0.95},
				time.Millisecond,
			)
		})
	}
	return nil
}
