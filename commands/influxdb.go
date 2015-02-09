package commands

import (
	"flag"
	"time"

	"github.com/ostrost/ostent/ostent"
	"github.com/ostrost/ostent/types"
	"github.com/rzab/go-metrics/influxdb"
)

type influx struct {
	logger      *Logger
	RefreshFlag types.PeriodValue
	ServerAddr  string // types.BindValue
	Database    string
	Username    string
	Password    string
}

func influxdbCommandLine(cli *flag.FlagSet) CommandLineHandler {
	ix := &influx{
		logger:      NewLogger("[ostent sendto-influxdb] "),
		RefreshFlag: types.PeriodValue{Duration: types.Duration(10 * time.Second)}, // 10s default
		// ServerAddr:  types.NewBindValue(8086),
	}
	cli.Var(&ix.RefreshFlag, "influxdb-refresh", "InfluxDB refresh interval")
	cli.StringVar(&ix.ServerAddr, "sendto-influxdb", "", "InfluxDB server address")
	cli.StringVar(&ix.Database, "influxdb-database", "ostent", "InfluxDB database")
	cli.StringVar(&ix.Username, "influxdb-username", "", "InfluxDB username")
	cli.StringVar(&ix.Password, "influxdb-password", "", "InfluxDB password")
	return func() (AtexitHandler, bool, error) {
		if ix.ServerAddr == "" {
			return nil, false, nil
		}
		ostent.AddBackground(func(defaultPeriod types.PeriodValue) {
			/* if ix.RefreshFlag.Duration == 0 { // if .RefreshFlag had no default
				ix.RefreshFlag = defaultPeriod
			} */
			go influxdb.Influxdb(ostent.Reg1s.Registry, time.Duration(ix.RefreshFlag.Duration), &influxdb.Config{
				Host:     ix.ServerAddr, //.String(),
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
