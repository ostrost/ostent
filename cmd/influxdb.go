package cmd

import (
	influxdb "github.com/vrischmann/go-metrics-influxdb"

	"github.com/ostrost/ostent/ostent"
	"github.com/ostrost/ostent/params"
)

func InfluxRun(ix params.InfluxParams) error {
	if ix.ServerAddr.String() != "" {
		u := ix.URL     // copy
		u.RawQuery = "" // reset query string
		ostent.AddBackground(func() {
			ostent.AllExporters.AddExporter("influxdb")
			go influxdb.InfluxDB(ostent.Reg1s.Registry,
				ix.Delay.D,
				u.String(),
				ix.Database,
				ix.Username,
				ix.Password,
			)
		})
	}
	return nil
}
