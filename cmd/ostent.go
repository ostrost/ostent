package cmd

import (
	"os"
	"time"

	"github.com/spf13/cobra"
	// "github.com/spf13/viper"

	"github.com/ostrost/ostent/cmd/cmdcobra"
	"github.com/ostrost/ostent/flags"
	"github.com/ostrost/ostent/ostent"
)

// OstentCmd represents the base command when called without any subcommands
var OstentCmd = &cobra.Command{
	SilenceUsage: true,
	Use:          "ostent",
	Short:        "Ostent is a metrics tool",
	Long: `Ostent collects system metrics and put them on display.
Optionally exports them to metrics servers.

Specify --influxdb-url  to enable exporting to InfluxDB.
Specify --graphite-host to enable exporting to Graphite.
Specify --librato-email and --librato-token  to enable exporting to Librato.
`,

	PostRunE: cmdcobra.PostRuns.RunE,
	PreRunE:  cmdcobra.PreRuns.RunE,
	// RunE in main.{bin,dev}.go
}

var DelayFlags = flags.DelayBounds{
	Max: flags.Delay{Duration: 10 * time.Minute},
	Min: flags.Delay{Duration: time.Second},
	// 10m and 1s are corresponding defaults
}

// Execute adds all child commands to the ostent command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the ostentCmd.
func Execute() {
	if err := OstentCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}

func init() {
	// cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.

	// OstentCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ostent.yaml)")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.

	OstentCmd.PersistentFlags().VarP(&OstentBind, "bind", "b", "Bind `address`")
	OstentCmd.PersistentFlags().Var(&DelayFlags.Max, "max-delay", "Collect and display maximum `delay`")
	OstentCmd.PersistentFlags().VarP(&DelayFlags.Min, "min-delay", "d", "Collect and display minimum `delay`")
	OstentCmd.PersistentFlags().BoolVarP(&Vflag, "version", "v", false, "Display version and exit")

	gr := &Graphite{
		DelayFlag:  flags.Delay{Duration: 10 * time.Second}, // 10s default
		ServerAddr: flags.NewBind(2003),
	}
	OstentCmd.PersistentFlags().Var(&gr.DelayFlag, "graphite-delay", "Graphite `delay`")
	OstentCmd.PersistentFlags().Var(&gr.ServerAddr, "graphite-host", "Graphite `host`")

	ix := &Influx{
		DelayFlag: flags.Delay{Duration: 10 * time.Second}, // 10s default
	}
	OstentCmd.PersistentFlags().Var(&ix.DelayFlag, "influxdb-delay", "InfluxDB `delay`")
	OstentCmd.PersistentFlags().StringVar(&ix.URL, "influxdb-url", "", "InfluxDB server `URL`")
	OstentCmd.PersistentFlags().StringVar(&ix.Database, "influxdb-database", "ostent", "InfluxDB `database`")
	OstentCmd.PersistentFlags().StringVar(&ix.Username, "influxdb-username", "", "InfluxDB `username`")
	OstentCmd.PersistentFlags().StringVar(&ix.Password, "influxdb-password", "", "InfluxDB `password`")

	hostname, _ := ostent.GetHN()
	lr := &Librato{
		DelayFlag: flags.Delay{Duration: 10 * time.Second}, // 10s default
	}
	OstentCmd.PersistentFlags().Var(&lr.DelayFlag, "librato-delay", "Librato `delay`")
	OstentCmd.PersistentFlags().StringVar(&lr.Email, "librato-email", "", "Librato `email`")
	OstentCmd.PersistentFlags().StringVar(&lr.Token, "librato-token", "", "Librato `token`")
	OstentCmd.PersistentFlags().StringVar(&lr.Source, "librato-source", hostname, "Librato `source`")

	cmdcobra.PreRuns.Adds(FixDelayFlags,
		OstentVersionRun, // version goes into PreRuns first
		gr.Run,
		ix.Run,
		lr.Run)
}

func FixDelayFlags() error {
	if DelayFlags.Max.Duration < DelayFlags.Min.Duration {
		DelayFlags.Max.Duration = DelayFlags.Min.Duration
	}
	return nil
}

/*
// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName(".ostent") // name of config file (without extension)
	viper.AddConfigPath("$HOME")  // adding home directory as first search path
	viper.AutomaticEnv()          // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
*/

var OstentBind = flags.NewBind(8050)
