package cmd

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/influxdata/toml"
	"github.com/spf13/cobra"

	"github.com/ostrost/ostent/flags"
	"github.com/ostrost/ostent/internal"
	"github.com/ostrost/ostent/internal/agent"
	"github.com/ostrost/ostent/internal/config"
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
			log.Printf("Config overview:\n%s", text)
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

func printableConfigText(text []byte) []byte {
	lines := bytes.Split(text, []byte("\n"))
	for _, replaceString := range [][2]string{
		{"password = ", `"********"`},
		{"api_token = ", `"****************"`},
	} {
		replace := [2][]byte{[]byte(replaceString[0]), []byte(replaceString[1])}
		for i := range lines {
			if j := bytes.Index(lines[i], replace[0]); j != -1 {
				lines[i] = append(append(lines[i][:j], replace[0]...), replace[1]...)
			}
		}
	}
	return bytes.Join(lines, []byte("\n"))
}

func printableConfig(rconfig *config.Config) ([]byte, error) {
	type agentSelect struct {
		Debug         bool `toml:",omitempty"`
		Quiet         bool `toml:",omitempty"`
		FlushInterval internal.Duration
		Interval      internal.Duration
		RoundInterval bool `toml:",omitempty"`
	}
	rt := struct {
		AgentSelect *agentSelect
		Ins         map[string]map[string]interface{}
		Outs        map[string]map[string]interface{}
	}{
		AgentSelect: &agentSelect{
			Debug:         rconfig.Agent.Debug,
			Quiet:         rconfig.Agent.Quiet,
			FlushInterval: rconfig.Agent.FlushInterval,
			Interval:      rconfig.Agent.Interval,
			RoundInterval: rconfig.Agent.RoundInterval,
		},
		Ins:  make(map[string]map[string]interface{}),
		Outs: make(map[string]map[string]interface{}),
	}
	for _, input := range rconfig.Inputs {
		rt.Ins["enable_"+input.Name] = nil
	}
	for _, output := range rconfig.Outputs {
		rt.Outs["enable_"+output.Name] = nil
	}
	text, err := toml.Marshal(rt)
	if err != nil {
		return []byte{}, err
	}
	return printableConfigText(text), nil
}
