package cmd

import (
	"github.com/ostrost/ostent/ostent"
	"github.com/ostrost/ostent/params"
)

type Influxdb struct {
	Namedrop namedrop
	URLs     []string `toml:"urls"`
	Username string
	Password string
	Database string
}

func InfluxRun(elisting *ostent.ExportingListing, tabs *tables, iends params.InfluxEndpoints) error {
	for _, value := range iends.Values {
		if value.ServerAddr.String() != "" {
			elisting.AddExporter("InfluxDB", value)
			u := value.URL  // copy
			u.RawQuery = "" // reset query string
			tabs.add(struct {
				Outputs []Influxdb `toml:"outputs.influxdb"`
			}{
				Outputs: []Influxdb{{
					Namedrop: commonNamedrop,
					URLs:     []string{u.String()},
					Username: value.Username,
					Password: value.Password,
					Database: value.Database,

					// TODO value.Tags is ignored
				}}})
		}
	}
	return nil
}
