package cmd

import (
	"github.com/ostrost/ostent/internal/config"
	"github.com/ostrost/ostent/ostent"
	"github.com/ostrost/ostent/params"
)

type Graphite struct {
	Namedrop namedrop
	Servers  []string
	Prefix   string // hard-coded
}

// type ConfigType

func GraphiteRun(elisting *ostent.ExportingListing, cconfig *config.Config, gends params.GraphiteEndpoints) error {
	for _, value := range gends.Values {
		if value.ServerAddr.Host != "" {
			elisting.AddExporter("Graphite", value)
			if err := cconfig.LoadInterface("/internal/graphite/config", struct {
				Outputs []Graphite `toml:"outputs.graphite"`
			}{
				Outputs: []Graphite{{
					Namedrop: commonNamedrop,
					Servers:  []string{value.ServerAddr.String()},
					Prefix:   "ostent", // hard-coded
				}}}); err != nil {
				return err
			}
		}
	}
	return nil
}
