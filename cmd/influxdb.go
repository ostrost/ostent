package cmd

import (
	influxdb "github.com/vrischmann/go-metrics-influxdb"

	"github.com/ostrost/ostent/ostent"
	"github.com/ostrost/ostent/params"
)

func InfluxRun(iends params.InfluxEndpoints) error {
	for _, value := range iends.Values {
		if value.ServerAddr.String() != "" {
			ostent.AddBackground(InfluxRunFunc(value))
		}
	}
	return nil
}

func InfluxRunFunc(value params.InfluxEndpoint) func() {
	return func() {
		u := value.URL  // copy
		u.RawQuery = "" // reset query string
		ostent.AllExporters.AddExporter("influxdb")
		go influxdb.InfluxDB(ostent.Reg1s.Registry,
			value.Delay.D,
			u.String(),
			value.Database,
			value.Username,
			value.Password,
		)
	}
}
