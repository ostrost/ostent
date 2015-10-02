package commands

import (
	"flag"
	"time"

	influxdb "github.com/vrischmann/go-metrics-influxdb"

	"github.com/ostrost/ostent/commands/extpoints"
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

func (_ Influxes) SetupFlagSet(cli *flag.FlagSet) extpoints.CommandLineHandler {
	ix := &Influx{
		DelayFlag: flags.Delay{Duration: 10 * time.Second}, // 10s default
	}
	cli.Var(&ix.DelayFlag, "influxdb-delay", "InfluxDB `delay`")
	cli.StringVar(&ix.URL, "influxdb-url", "", "InfluxDB server `URL`")
	cli.StringVar(&ix.Database, "influxdb-database", "ostent", "InfluxDB `database`")
	cli.StringVar(&ix.Username, "influxdb-username", "", "InfluxDB `username`")
	cli.StringVar(&ix.Password, "influxdb-password", "", "InfluxDB `password`")
	return func() (extpoints.AtexitHandler, bool, error) {
		if ix.URL == "" {
			return nil, false, nil
		}
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
		return nil, false, nil
	}
}

type Influxes struct{}

func init() {
	extpoints.CommandLines.Register(Influxes{}, "influxdb")
}
