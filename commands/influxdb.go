package commands

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/ostrost/ostent"
	"github.com/ostrost/ostent/types"
	"github.com/rcrowley/go-metrics/influxdb"
)

type influx struct {
	logger     *loggerWriter
	ServerAddr types.BindValue
	Database   string
	Username   string
	Password   string
}

func influxdbCommandLine(cli *flag.FlagSet) commandLineHandler {
	ix := &influx{
		logger: &loggerWriter{
			log.New(os.Stderr, "[ostent sendto-influxdb] ", log.LstdFlags),
		},
		ServerAddr: types.NewBindValue(8086),
	}
	cli.Var(&ix.ServerAddr, "sendto-influxdb", "InfluxDB server address")
	cli.StringVar(&ix.Database, "influxdb-database", "ostent", "InfluxDB database")
	cli.StringVar(&ix.Username, "influxdb-username", "", "InfluxDB username")
	cli.StringVar(&ix.Password, "influxdb-password", "", "InfluxDB password")
	return func() (atexitHandler, bool, error) {
		if ix.ServerAddr.Host == "" {
			return nil, false, nil
		}
		ostent.AddBackground(func() {
			go influxdb.Influxdb(ostent.Reg1s.Registry, 1*time.Second, &influxdb.Config{
				Host:     ix.ServerAddr.String(),
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
