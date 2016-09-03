package cmd

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sort"
	"strings"
	"sync"
	"syscall"
	"unicode"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/influxdata/toml"

	"github.com/ostrost/ostent/flags"
	"github.com/ostrost/ostent/internal/agent"
	"github.com/ostrost/ostent/internal/config"
	"github.com/ostrost/ostent/ostent"
	"github.com/ostrost/ostent/params"

	// plugging outputs:
	_ "github.com/influxdata/telegraf/plugins/outputs/graphite"
	_ "github.com/influxdata/telegraf/plugins/outputs/influxdb"
	_ "github.com/influxdata/telegraf/plugins/outputs/librato"

	_ "github.com/ostrost/ostent/internal/plugins/outputs/ostent" // "ostent" output

	// plugging inputs:
	_ "github.com/influxdata/telegraf/plugins/inputs/system" // "{cpu,disk,mem,swap}" inputs

	_ "github.com/ostrost/ostent/procstat_ostent" // "procstat_ostent" input
	_ "github.com/ostrost/ostent/system_ostent"   // "{net,system}_ostent" inputs
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

// Bind is flag value.
var Bind = flags.NewBind("", 8050)

func initFlags() {
	rconfig := defaultConfig()
	tabs := &tables{}

	persistentPreRuns.add(versionRun)
	RootCmd.PersistentFlags().BoolVar(&versionFlag, "version", false, "Print version and exit")

	RootCmd.Flags().VarP(&Bind, "bind", "b", "Bind `address`")

	defaultInterval := config.NewConfig().Agent.Interval.Duration.String()
	intervals := struct{ Agent, CPU, Disk, Mem, NetOstent, ProcstatOstent, Swap string }{}
	for _, v := range []struct {
		pointer *string
		name    string
		usage   string
	}{
		{&intervals.CPU, "input-interval-cpu", "Interval for input: cpu"},
		{&intervals.Disk, "input-interval-disk", "Interval for input: disk"},
		{&intervals.Mem, "input-interval-mem", "Interval for input: mem"},
		{&intervals.NetOstent, "input-interval-net-ostent", "Interval for input: net_ostent"},
		{&intervals.ProcstatOstent, "input-interval-procstat-ostent", "Interval for input: procstat_ostent"},
		{&intervals.Swap, "input-interval-swap", "Interval for input: swap"},
		// ...
	} {
		RootCmd.Flags().StringVar(v.pointer, v.name, defaultInterval, v.usage)
	}
	RootCmd.Flags().StringVarP(&intervals.Agent, "interval", "d", defaultInterval, "Interval for agent and inputs")

	preRuns.add(func() error {

		if intervals.Agent != defaultInterval {
			interval := new(string)
			*interval = intervals.Agent
			rconfig.Agent.Interval, rconfig.Agent.FlushInterval = interval, interval
			rconfig.Inputs.System_ostent.Interval = interval // otherwise it's 1s per System_ostent default
		}

		for _, v := range []struct {
			pointer **string
			value   string
		}{
			{&rconfig.Inputs.CPU.Interval, intervals.CPU},
			{&rconfig.Inputs.Disk.Interval, intervals.Disk},
			{&rconfig.Inputs.Mem.Interval, intervals.Mem},
			{&rconfig.Inputs.Net_ostent.Interval, intervals.NetOstent},
			{&rconfig.Inputs.Procstat_ostent.Interval, intervals.ProcstatOstent},
			{&rconfig.Inputs.Swap.Interval, intervals.Swap},
		} {
			if v.value != defaultInterval {
				*v.pointer = new(string)
				**v.pointer = v.value
			}
		}
		tabs.add(rconfig)
		return nil
	})
	var elisting ostent.ExportingListing

	if gends := params.NewGraphiteEndpoints(flags.NewBind("127.0.0.1", 2003)); true {
		preRuns.add(func() error { return graphiteRun(&elisting, tabs, gends) })
		RootCmd.Flags().Var(&gends, "graphite", "Graphite exporting `endpoint(s)`")
	}

	if iends := params.NewInfluxEndpoints("ostent"); true {
		preRuns.add(func() error { return influxRun(&elisting, tabs, iends) })
		RootCmd.Flags().Var(&iends, "influxdb", "InfluxDB exporting `endpoint(s)`")
		RootCmd.Example += "InfluxDB params:\n" + paramsUsage(func(f *pflag.FlagSet) {
			param := &iends.Default // shortcut, f does not alter it
			// f.Var(&param.ServerAddr, "0", "InfluxDB server `address`")
			f.StringVar(&param.Database, "2", param.Database, "InfluxDB `database`")
			f.StringVar(&param.Username, "3", param.Username, "InfluxDB `username`")
			f.StringVar(&param.Password, "4", param.Password, "InfluxDB `password`")
		}) + "  Any extra parameters become tags in every metrics post to InfluxDB server.\n"
	}

	hostname, _ := os.Hostname()
	if lends := params.NewLibratoEndpoints(hostname); true {
		preRuns.add(func() error { return libratoRun(&elisting, tabs, lends) })
		RootCmd.Flags().Var(&lends, "librato", "Librato exporting `parameter(s)`")
		RootCmd.Example += "Librato params:\n" + paramsUsage(func(f *pflag.FlagSet) {
			param := &lends.Default // shortcut, f does not alter it
			f.StringVar(&param.Source, "2", param.Source, "Librato `source`")
			f.StringVar(&param.Email, "3", param.Email, "Librato `email`")
			f.StringVar(&param.Token, "4", param.Token, "Librato `token`")
		})
	}
	RootCmd.Example = strings.TrimRight(RootCmd.Example, "\n")
	preRuns.add(func() error {
		ostent.Exporting = elisting.ExportingList
		sort.Stable(ostent.Exporting)
		return nil
	})

	// /*
	ostent.AddBackground(func() {
		if err := mainAgent(tabs); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}) // */
}

// paramsUsage returns formatted usage of a FlagSet set by setf.
// All flags assumed to be params thus formatting trims dashes.
// The flag names supposed to be digits so it strips them likewise.
func paramsUsage(setf func(*pflag.FlagSet)) string {
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

func defaultConfig() oneConfig {
	var (
		oneSecond = "1s"
		ignoreFs  = []string{"tmpfs", "devtmpfs"}
	)
	var (
		on  = func() *inputConfig { return &inputConfig{} }
		ion = func() *[]struct{} { return &[]struct{}{struct{}{}} }
	)

	return oneConfig{
		Outputs: outputs{Ostent: ion()},
		Inputs: inputs{
			CPU:             on(),
			Disk:            &diskInput{Ignore_fs: &ignoreFs},
			Mem:             on(),
			Net_ostent:      on(),
			Procstat_ostent: on(),
			Swap:            on(),
			System_ostent:   &inputConfig{&oneSecond},
		}}
}

type oneConfig struct {
	Agent   agentConfig
	Outputs outputs
	Inputs  inputs
}

type agentConfig struct {
	Interval      *string `toml:",omitempty" yaml:",omitempty"`
	FlushInterval *string `toml:",omitempty" yaml:",omitempty"`
}

type inputConfig struct {
	Interval *string `toml:",omitempty" yaml:",omitempty"`
}

type diskInput struct {
	Interval  *string   `toml:",omitempty" yaml:",omitempty"` // common inputConfig
	Ignore_fs *[]string `toml:",omitempty" yaml:",omitempty"`
}

type outputs struct {
	Ostent   *[]struct{} `toml:",omitempty" yaml:",omitempty"`
	Influxdb *[]struct {
		Username, Password, Database string
		Namedrop, URLs               []string
	} `toml:",omitempty" yaml:",omitempty"`
}

type inputs struct {
	CPU             *inputConfig `toml:",omitempty" yaml:",omitempty"`
	Disk            *diskInput   `toml:",omitempty" yaml:",omitempty"`
	Mem             *inputConfig `toml:",omitempty" yaml:",omitempty"`
	Net_ostent      *inputConfig `toml:",omitempty" yaml:",omitempty"`
	Procstat_ostent *inputConfig `toml:",omitempty" yaml:",omitempty"`
	Swap            *inputConfig `toml:",omitempty" yaml:",omitempty"`
	System_ostent   *inputConfig `toml:",omitempty" yaml:",omitempty"`
}

type namedrop []string

var commonNamedrop = namedrop{
	"system_ostent",
	"procstat",
	"procstat_ostent",
}

func mainAgent(tabs *tables) error {
	reload := make(chan bool, 1)
	reload <- true
	for <-reload {
		reload <- false

		c := config.NewConfig()

		configText, err := tabs.marshal()
		if err != nil {
			return err
		}

		configPath := "/runtime/config"
		log.Printf("#%s.toml:\n%s", configPath, printableConfigText(configText))

		configTab, err := config.ParseContents(configText)
		if err != nil {
			return err
		}
		if err := c.LoadTable(configPath, configTab); err != nil {
			return err
		}

		c.Agent.Quiet = true

		ag, err := agent.NewAgent(c)
		if err != nil {
			return err
		}

		if err := ag.Connect(); err != nil {
			return err
		}

		shutdown := make(chan struct{})
		signals := make(chan os.Signal)
		signal.Notify(signals, os.Interrupt, syscall.SIGHUP)
		go func() {
			sig := <-signals
			if sig == os.Interrupt {
				close(shutdown)
			}
			if sig == syscall.SIGHUP {
				log.Printf("Reloading config\n")
				<-reload
				reload <- true
				close(shutdown)
			}
		}()

		if err := ag.Run(shutdown); err != nil {
			return err
		}
	}
	return nil
}

type tables struct{ list []interface{} }

func (tabs *tables) add(in interface{}) { tabs.list = append(tabs.list, in) }

func (tabs tables) marshal() ([]byte, error) {
	text := []byte{}
	for _, in := range tabs.list {
		add, err := toml.Marshal(in)
		if err != nil {
			return text, err
		}
		text = append(text, add...)
	}
	return text, nil
}

func printableConfigText(text []byte) string {
	lines := strings.Split(string(text), "\n")
	for _, replace := range [][2]string{
		{"password = ", `"********"`},
		{"api_token = ", `"****************"`},
	} {
		for i := range lines {
			if j := strings.Index(lines[i], replace[0]); j != -1 {
				lines[i] = lines[i][:j] + replace[0] + replace[1]
			}
		}
	}
	return strings.Join(lines, "\n")
}
