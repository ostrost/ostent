package cmd

import (
	influxdb "github.com/vrischmann/go-metrics-influxdb"

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

func InfluxRun(elisting *ostent.ExportingListing, cconfig configer, iends params.InfluxEndpoints) error {
	for _, value := range iends.Values {
		if value.ServerAddr.String() != "" {
			elisting.AddExporter("InfluxDB", value)
			u := value.URL  // copy
			u.RawQuery = "" // reset query string
			err := cconfig.LoadInterface("/internal/influxdb/config", struct {
				Outputs []Influxdb `toml:"outputs.influxdb"`
			}{
				Outputs: []Influxdb{{
					Namedrop: commonNamedrop,
					URLs:     []string{u.String()},
					Username: value.Username,
					Password: value.Password,
					Database: value.Database,

					// TODO value.Tags is ignored
					// TODO value.Delay becomes meaningless
				}}})
			if err != nil {
				return err
			}
			ostent.AddBackground(InfluxRunFunc(value))
		}
	}
	return nil
}

func InfluxRunFunc(value params.InfluxEndpoint) func() {
	return func() {
		u := value.URL  // copy
		u.RawQuery = "" // reset query string
		go influxdb.InfluxDBWithTags(ostent.Reg1s.Registry,
			value.Delay.D,
			u.String(),
			value.Database,
			value.Username,
			value.Password,
			value.Tags,
		)
	}
}
