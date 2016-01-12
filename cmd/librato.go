package cmd

import (
	"time"

	librato "github.com/mihasya/go-metrics-librato"

	"github.com/ostrost/ostent/ostent"
	"github.com/ostrost/ostent/params"
)

func LibratoRun(lr params.LibratoParams) error {
	if lr.Email != "" {
		ostent.AddBackground(func() {
			ostent.AllExporters.AddExporter("librato")
			go librato.Librato(ostent.Reg1s.Registry,
				lr.Delay.D,
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
