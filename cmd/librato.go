package cmd

import (
	"time"

	librato "github.com/mihasya/go-metrics-librato"

	"github.com/ostrost/ostent/ostent"
	"github.com/ostrost/ostent/params"
)

func LibratoRun(elisting *ostent.ExportingListing, lends params.LibratoEndpoints) error {
	for _, value := range lends.Values {
		if value.Email != "" {
			elisting.AddExporter("Librato", value)
			ostent.AddBackground(LibratoRunFunc(value))
		}
	}
	return nil
}

func LibratoRunFunc(value params.LibratoEndpoint) func() {
	return func() {
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
