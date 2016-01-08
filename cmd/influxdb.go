package cmd

import (
	influxdb "github.com/vrischmann/go-metrics-influxdb"

	"github.com/ostrost/ostent/flags"
	"github.com/ostrost/ostent/ostent"
)

type Influx struct {
	DelayFlag flags.Delay
	URL       string
	Database  string
	Username  string
	Password  string
}

func (ix *Influx) Run() error {
	if ix.URL != "" {
		ostent.AddBackground(func() {
			ostent.AllExporters.AddExporter("influxdb")
			go influxdb.InfluxDB(ostent.Reg1s.Registry,
				ix.DelayFlag.Duration,
				ix.URL,
				ix.Database,
				ix.Username,
				ix.Password,
			)
		})
	}
	return nil
}
