package cmd

import (
	"os"
	"sort"
	"strings"
	"sync"
	"time"
	"unicode"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	// "github.com/spf13/viper"

	"github.com/ostrost/ostent/flags"
	"github.com/ostrost/ostent/ostent"
	"github.com/ostrost/ostent/params"
)

var (
	persistentPostRuns runs // a list of funcs to be cobra.Command's PersistentPostRunE.
	persistentPreRuns  runs // a list of funcs to be cobra.Command's PersistentPreRunEE.
	preRuns            runs // a list of funcs to be cobra.Command's PreRunE.
)

type runs struct {
	mutex sync.Mutex     // protect everything e.g. list
	list  []func() error // the list to have
}

func (rs *runs) add(f func() error) {
	rs.mutex.Lock()
	defer rs.mutex.Unlock()
	rs.list = append(rs.list, f)
}

func (rs *runs) runE(*cobra.Command, []string) error {
	rs.mutex.Lock()
	defer rs.mutex.Unlock()
	for _, run := range rs.list {
		if err := run(); err != nil {
			return err
		}
	}
	return nil
}

var (
	// DelayFlags sets min and max for any delay.
	DelayFlags = flags.DelayBounds{
		Max: flags.Delay{Duration: 10 * time.Minute},
		Min: flags.Delay{Duration: time.Second},
		// 10m and 1s are corresponding defaults
	}

	// OstentBind is the flag value.
	OstentBind = flags.NewBind("", 8050)

	// OstentCmd represents the base command when called without any subcommands
	OstentCmd = &cobra.Command{
		SilenceUsage: true,
		Use:          "ostent",
		Short:        "Ostent is a metrics tool.",
		Long: `Ostent collects system metrics and put them on display.
Optionally exports them to metrics servers.

To continuously export collected metrics use graphite, influxdb and/or librato flags.
Use multiple flags and/or use comma separated endpoints for the same kind.`,
		Example: `ostent --graphite 10.0.0.1,10.0.0.2:2004\?delay=30s
ostent --influxdb http://10.0.0.3\?delay=60s
ostent --librato \?email=EMAIL\&token=TOKEN

`,

		PersistentPostRunE: persistentPostRuns.runE,
		PersistentPreRunE:  persistentPreRuns.runE,
		PreRunE:            preRuns.runE,
	}
)

func init() {
	// cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.

	// OstentCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ostent.yaml)")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.

	persistentPreRuns.add(OstentVersionRun)
	OstentCmd.PersistentFlags().BoolVar(&VersionFlag, "version", false, "Print version and exit")

	OstentCmd.Flags().VarP(&OstentBind, "bind", "b", "Bind `address`")
	OstentCmd.Flags().Var(&DelayFlags.Max, "max-delay", "Maximum for display `delay`")
	OstentCmd.Flags().VarP(&DelayFlags.Min, "min-delay", "d", "Collection and display minimum `delay`")

	preRuns.add(func() error {
		if DelayFlags.Max.Duration < DelayFlags.Min.Duration {
			DelayFlags.Max.Duration = DelayFlags.Min.Duration
		}
		return nil
	})

	var elisting ostent.ExportingListing

	if gends := params.NewGraphiteEndpoints(10*time.Second, flags.NewBind("127.0.0.1", 2003)); true {
		preRuns.add(func() error { return GraphiteRun(&elisting, gends) })
		OstentCmd.Flags().Var(&gends, "graphite", "Graphite exporting `endpoint(s)`")
		OstentCmd.Example += "Graphite params:\n" + ParamsUsage(func(f *pflag.FlagSet) {
			param := &gends.Default // shortcut, f does not alter it
			// f.Var(&param.ServerAddr, "0", "Graphite server `host[:port]`")
			f.Var(&param.Delay, "1", "Graphite exporting `delay`")
		})
	}

	if iends := params.NewInfluxEndpoints(10*time.Second, "ostent"); true {
		preRuns.add(func() error { return InfluxRun(&elisting, iends) })
		OstentCmd.Flags().Var(&iends, "influxdb", "InfluxDB exporting `endpoint(s)`")
		OstentCmd.Example += "InfluxDB params:\n" + ParamsUsage(func(f *pflag.FlagSet) {
			param := &iends.Default // shortcut, f does not alter it
			// f.Var(&param.ServerAddr, "0", "InfluxDB server `address`")
			f.Var(&param.Delay, "1", "InfluxDB exporting `delay`")
			f.StringVar(&param.Database, "2", param.Database, "InfluxDB `database`")
			f.StringVar(&param.Username, "3", param.Username, "InfluxDB `username`")
			f.StringVar(&param.Password, "4", param.Password, "InfluxDB `password`")
		}) + "  Any extra parameters become tags in every metrics post to InfluxDB server.\n"
	}

	hostname, _ := os.Hostname()
	if lends := params.NewLibratoEndpoints(10*time.Second, hostname); true {
		preRuns.add(func() error { return LibratoRun(&elisting, lends) })
		OstentCmd.Flags().Var(&lends, "librato", "Librato exporting `parameter(s)`")
		OstentCmd.Example += "Librato params:\n" + ParamsUsage(func(f *pflag.FlagSet) {
			param := &lends.Default // shortcut, f does not alter it
			f.Var(&param.Delay, "1", "Librato exporting `delay`")
			f.StringVar(&param.Source, "2", param.Source, "Librato `source`")
			f.StringVar(&param.Email, "3", param.Email, "Librato `email`")
			f.StringVar(&param.Token, "4", param.Token, "Librato `token`")
		})
	}
	OstentCmd.Example = strings.TrimRight(OstentCmd.Example, "\n")
	preRuns.add(func() error {
		ostent.Exporting = elisting.ExportingList
		sort.Stable(ostent.Exporting)
		return nil
	})
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

// ParamsUsage returns formatted usage of a FlagSet set by setf.
// All flags assumed to be params thus formatting trims dashes.
// The flag names supposed to be digits so it strips them likewise.
func ParamsUsage(setf func(*pflag.FlagSet)) string {
	cmd := cobra.Command{}
	setf(cmd.Flags())
	lines := strings.Split(cmd.NonInheritedFlags().FlagUsages(), "\n")
	for i := range lines {
		lines[i] = strings.TrimPrefix(lines[i], "      --")
		lines[i] = strings.TrimLeftFunc(lines[i], unicode.IsDigit)
		if lines[i] != "" {
			lines[i] = " " + lines[i]
		}
	}
	return strings.Join(lines, "\n")
}
