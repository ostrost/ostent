package cmd

import (
	"time"

	librato "github.com/mihasya/go-metrics-librato"

	"github.com/ostrost/ostent/ostent"
	"github.com/ostrost/ostent/params"
)

func LibratoRun(lends params.LibratoEndpoints) error {
	for _, value := range lends.Values {
		if value.Email != "" {
			ostent.AddBackground(LibratoRunFunc(value))
		}
	}
	return nil
}

func LibratoRunFunc(value params.LibratoEndpoint) func() {
	return func() {
		ostent.AllExporters.AddExporter("librato")
		go librato.Librato(ostent.Reg1s.Registry,
			value.Delay.D,
			value.Email,
			value.Token,
			value.Source,
			[]float64{0.95},
			time.Millisecond,
		)
	}
}
