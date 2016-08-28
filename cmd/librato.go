package cmd

import (
	"github.com/ostrost/ostent/ostent"
	"github.com/ostrost/ostent/params"
)

type Librato struct {
	Namedrop  namedrop
	ApiUser   string
	ApiToken  string
	SourceTag string
	Template  string // hard-coded
}

func LibratoRun(elisting *ostent.ExportingListing, tabs *tables, lends params.LibratoEndpoints) error {
	for _, value := range lends.Values {
		if value.Email != "" {
			elisting.AddExporter("Librato", value)
			tabs.add(struct {
				Outputs []Librato `toml:"outputs.librato"`
			}{
				Outputs: []Librato{{
					Namedrop:  commonNamedrop,
					ApiUser:   value.Email,
					ApiToken:  value.Token,
					SourceTag: value.Source,
					Template:  "host.tags.measurement.field", // hard-coded
				}}})
		}
	}
	return nil
}
