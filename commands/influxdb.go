package commands

import (
	"flag"
	"time"

	influxdb "github.com/ostrost/ostent/commands/go-metrics-influxdb"
	"github.com/ostrost/ostent/flags"
	"github.com/ostrost/ostent/ostent"
)

type Influx struct {
	// Logger      *Logger
	DelayFlag flags.Delay
	URL       string
	Database  string
	Username  string
	Password  string
}

func influxdbCommandLine(cli *flag.FlagSet) CommandLineHandler {
	ix := &Influx{
		// Logger:      NewLogger("[ostent sendto-influxdb] "),
		DelayFlag: flags.Delay{Duration: 10 * time.Second}, // 10s default
	}
	cli.Var(&ix.DelayFlag, "influxdb-delay", "InfluxDB `delay`")
	cli.StringVar(&ix.URL, "influxdb-url", "", "InfluxDB server `URL`")
	cli.StringVar(&ix.Database, "influxdb-database", "ostent", "InfluxDB `database`")
	cli.StringVar(&ix.Username, "influxdb-username", "", "InfluxDB `username`")
	cli.StringVar(&ix.Password, "influxdb-password", "", "InfluxDB `password`")
	return func() (AtexitHandler, bool, error) {
		if ix.URL == "" {
			return nil, false, nil
		}
		ostent.AddBackground(func() {
			go influxdb.Influxdb(ostent.Reg1s.Registry, ix.DelayFlag.Duration, &influxdb.Config{
				URL:      ix.URL,
				Database: ix.Database,
				Username: ix.Username,
				Password: ix.Password,
			})
		})
		return nil, false, nil
	}
}

func init() {
	AddCommandLine(influxdbCommandLine)
}
