package cmd

import (
	"time"

	librato "github.com/mihasya/go-metrics-librato"

	"github.com/ostrost/ostent/internal/config"
	"github.com/ostrost/ostent/ostent"
	"github.com/ostrost/ostent/params"
)

type Librato struct {
	ApiUser   string
	ApiToken  string
	SourceTag string
	Template  string // hard-coded

	Namedrop []string // general output preference
}

func LibratoRun(elisting *ostent.ExportingListing, cconfig *config.Config, lends params.LibratoEndpoints) error {
	for _, value := range lends.Values {
		if value.Email != "" {
			elisting.AddExporter("Librato", value)
			err := cconfig.LoadInterface("/internal/librato/config", struct {
				Outputs []Librato `toml:"outputs.librato"`
			}{
				Outputs: []Librato{{
					ApiUser:   value.Email,
					ApiToken:  value.Token,
					SourceTag: value.Source,
					Template:  "host.tags.measurement.field", // hard-coded

					// TODO value.Delay becomes meaningless
					Namedrop: []string{"system_ostent"},
				}}})
			if err != nil {
				return err
			}
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
