package cmd

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/influxdata/toml"
	"github.com/spf13/cobra"

	"github.com/ostrost/ostent/flags"
	"github.com/ostrost/ostent/internal"
	"github.com/ostrost/ostent/internal/agent"
	"github.com/ostrost/ostent/internal/config"
	"github.com/ostrost/ostent/internal/models"
	"github.com/ostrost/ostent/ostent"

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
	rconfig := config.NewConfig()

	persistentPreRuns.add(versionRun)
	RootCmd.PersistentFlags().BoolVar(&versionFlag, "version", false, "Print version and exit")

	RootCmd.Flags().VarP(&Bind, "bind", "b", "Bind `address`")

	defaultInterval := rconfig.Agent.Interval.Duration.String()
	intervals := struct{ Agent string }{}
	/*
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
	*/
	RootCmd.Flags().StringVarP(&intervals.Agent, "interval", "d", defaultInterval, "Interval for agent and inputs")

	_ = defaultInterval
	/*
		preRuns.add(func() error {
			if intervals.Agent != defaultInterval {
				var interval internal.Duration
				_ = interval.UnmarshalTOML([]byte(fmt.Sprintf("%q", intervals.Agent)))
				// TODO set rconfig.Agent.{Flush,}Interval, rconfig.Outputs["system_ostent"].Interval
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
			return nil
		})
	*/

	preRuns.add(func() error {
		ostent.AddBackground(func() {
			if err := mainAgent(rconfig); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		})
		return nil
	})
}

func mainAgent(rconfig *config.Config) error {
	reload := make(chan bool, 1)
	reload <- true
	for <-reload {
		reload <- false

		c := rconfig // config.NewConfig() // TODO rconfig copy
		c.Agent.Quiet = true

		if tab, err := readConfig(rconfig); err != nil {
			return err
		} else if tab != nil {
			if err := c.LoadTable("/runtime/config", tab); err != nil {
				return err
			}
		}

		if text, err := printableConfig(rconfig); err != nil {
			return err
		} else {
			log.Printf("Runtime config:\n%s", text)
			/*
				var system_ostent *models.RunningInput
				for _, ri := range rconfig.Inputs {
					if ri.Name == "system_ostent" {
						system_ostent = ri
					}
				}
				log.Printf("ins...._system_ostent.interval = %q\n",
					system_ostent.Config.Interval) */
		}

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

func printableConfigText(text string) string {
	lines := strings.Split(text, "\n")
	for _, replace := range [][2]string{
		{"password = ", `"PASSWORD"`},
		{"api_token = ", `"API_TOKEN"`},
	} {
		for i := range lines {
			if j := strings.Index(lines[i], replace[0]); j != -1 {
				lines[i] = lines[i][:j] + replace[0] + replace[1]
			}
		}
	}
	return strings.Join(lines, "\n")
}

func printableConfig(rconfig *config.Config) (string, error) {
	agtext, err := toml.Marshal(struct{ Agent *config.AgentConfig }{rconfig.Agent})
	if err != nil {
		return "", err
	}

	text := string(agtext) + "[inputs]\n"
	for _, in := range rconfig.Inputs {
		subtext, err := printableTable("inputs", in.Name, in.Input,
			printableInput(in.Config), printableFilter(in.Config.Filter))
		if err != nil {
			return "", err
		}
		text += subtext
	}

	text += "[outputs]\n"
	for _, out := range rconfig.Outputs {
		subtext, err := printableTable("outputs", out.Name, out.Output,
			nil, printableFilter(out.Config.Filter))
		if err != nil {
			return "", err
		}
		text += subtext
	}

	return printableConfigText(text), nil
}

func printableTable(upname, name string, in1 interface{},
	in2 *printInput, in3 *printFilter) (string, error) {
	intext1, err := toml.Marshal(in1)
	if err != nil {
		return "", err
	}
	var intext2 []byte
	if in2 != nil {
		intext2, err = toml.Marshal(in2)
		if err != nil {
			return "", err
		}
	}
	var intext3 []byte
	if in3 != nil {
		intext3, err = toml.Marshal(in3)
		if err != nil {
			return "", err
		}
	}
	lines := strings.Split(string(intext1)+string(intext2)+string(intext3), "\n")
	for i := range lines {
		if lines[i] != "" {
			lines[i] = "        " + lines[i]
		}
	}
	return "    [" + upname + "." + name + "]\n" + strings.Join(lines, "\n"), nil
}

type printInput struct {
	NameOverride      string             `toml:",omitempty"`
	MeasurementPrefix string             `toml:",omitempty"`
	MeasurementSuffix string             `toml:",omitempty"`
	Tags              *map[string]string `toml:",omitempty"`
	Interval          internal.Duration  `toml:",omitempty"`
}

func printableInput(ic *models.InputConfig) *printInput {
	p := printInput{
		NameOverride:      ic.NameOverride,
		MeasurementPrefix: ic.MeasurementPrefix,
		MeasurementSuffix: ic.MeasurementSuffix,
		Interval:          internal.Duration{Duration: ic.Interval},
	}
	if ic.Tags != nil && len(ic.Tags) > 0 {
		p.Tags = &ic.Tags
	}
	if p == (printInput{}) {
		return nil
	}
	return &p
}

type printFilter struct {
	Filter struct {
		NameDrop   *[]string           `toml:",omitempty"`
		NamePass   *[]string           `toml:",omitempty"`
		FieldDrop  *[]string           `toml:",omitempty"`
		FieldPass  *[]string           `toml:",omitempty"`
		TagDrop    *[]models.TagFilter `toml:",omitempty"`
		TagPass    *[]models.TagFilter `toml:",omitempty"`
		TagExclude *[]string           `toml:",omitempty"`
		TagInclude *[]string           `toml:",omitempty"`
	}
}

func printableFilter(f models.Filter) *printFilter {
	if !f.IsActive {
		return nil
	}
	var p printFilter
	for _, pair := range []struct {
		in  *[]string
		out **[]string
	}{
		{&f.NameDrop, &p.Filter.NameDrop},
		{&f.NamePass, &p.Filter.NamePass},
		{&f.FieldDrop, &p.Filter.FieldDrop},
		{&f.FieldPass, &p.Filter.FieldPass},
		{&f.TagExclude, &p.Filter.TagExclude},
		{&f.TagInclude, &p.Filter.TagInclude},
	} {
		if len(*pair.in) > 0 {
			*pair.out = pair.in
		}
	}
	for _, pair := range []struct {
		in  *[]models.TagFilter
		out **[]models.TagFilter
	}{
		{&f.TagDrop, &p.Filter.TagDrop},
		{&f.TagPass, &p.Filter.TagPass},
	} {
		if len(*pair.in) > 0 {
			*pair.out = pair.in
		}
	}
	return &p
}
