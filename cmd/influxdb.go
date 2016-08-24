package cmd

import (
	"github.com/ostrost/ostent/internal/config"
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

func InfluxRun(elisting *ostent.ExportingListing, cconfig *config.Config, iends params.InfluxEndpoints) error {
	for _, value := range iends.Values {
		if value.ServerAddr.String() != "" {
			elisting.AddExporter("InfluxDB", value)
			u := value.URL  // copy
			u.RawQuery = "" // reset query string
			if err := cconfig.LoadInterface("/internal/influxdb/config", struct {
				Outputs []Influxdb `toml:"outputs.influxdb"`
			}{
				Outputs: []Influxdb{{
					Namedrop: commonNamedrop,
					URLs:     []string{u.String()},
					Username: value.Username,
					Password: value.Password,
					Database: value.Database,

					// TODO value.Tags is ignored
				}}}); err != nil {
				return err
			}
		}
	}
	return nil
}
