package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/influxdata/toml"
	"github.com/influxdata/toml/ast"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

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
	persistentPreRuns.add(versionRun)
	RootCmd.PersistentFlags().BoolVar(&versionFlag, "version", false, "Print version and exit")

	RootCmd.Flags().VarP(&Bind, "bind", "b", "Bind `address`")

	agentargs := &agentArguments{
		defaultInterval: config.NewConfig().Agent.Interval.Duration.String()}
	RootCmd.Flags().StringVar(&agentargs.interval, "interval",
		agentargs.defaultInterval, "Agent interval")

	preRuns.add(func() error {
		ostent.AddBackground(func() {
			if err := mainAgent(agentargs); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		})
		return nil
	})
}

type agentArguments struct {
	defaultInterval string // not an arg
	interval        string
}

func mainAgent(args *agentArguments) error {
	watchConfig()

	reload := make(chan bool, 1)
	reload <- true
	for <-reload {
		reload <- false

		var inputFilters []string
		var outputFilters []string

		c := config.NewConfig()
		c.OutputFilters = outputFilters
		c.InputFilters = inputFilters
		err := loadConfig(c, args)
		if err != nil {
			return err
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

func loadConfig(c *config.Config, args *agentArguments) error {
	c.Agent.Quiet = true // patch work

	if tab, err := readConfig(c); err != nil {
		return err
	} else if tab != nil {
		if err := c.LoadTable("/runtime/config", tab); err != nil {
			return err
		}
	}

	if args.interval != "" && args.interval != args.defaultInterval {
		var interval internal.Duration
		q := []byte(fmt.Sprintf("%q", args.interval))
		if err := interval.UnmarshalTOML(q); err != nil {
			return err
		}
		c.Agent.Interval = interval
	}

	if text, err := printableConfig(c); err != nil {
		return err
	} else {
		log.Printf("Effective runtime config:\n%s", text)
	}
	return nil
}

func printableConfigText(text string) string {
	lines := strings.Split(text, "\n")
	var newlines []string
rangelines:
	for i := range lines {
		for _, suffix := range []string{
			" = 0",
			` = ""`,
			" = []",
			` = "0s"`,
			" = false",
		} {
			if strings.HasSuffix(lines[i], suffix) {
				continue rangelines
			}
		}

		for _, replace := range [][2]string{
			{"password = ", `"PASSWORD"`},
			{"api_token = ", `"API_TOKEN"`},
		} {
			if j := strings.Index(lines[i], replace[0]); j != -1 {
				lines[i] = lines[i][:j] + replace[0] + replace[1]
			}
		}
		newlines = append(newlines, lines[i])
	}
	return strings.Join(newlines, "\n")
}

func printableConfig(rconfig *config.Config) (string, error) {
	agtext, err := toml.Marshal(struct{ Agent *config.AgentConfig }{rconfig.Agent})
	if err != nil {
		return "", err
	}

	text := string(agtext) + "[inputs]\n"
	for _, in := range rconfig.Inputs {
		tabtext, err := printableTable(in.Input, printableInput(in.Config),
			printableFilter(in.Config.Filter))
		if err != nil {
			return "", err
		}
		text += printableHeader("inputs", in.Name) + tabtext
	}

	text += "[outputs]\n"
	for _, out := range rconfig.Outputs {
		tabtext, err := printableTable(out.Output, nil,
			printableFilter(out.Config.Filter))
		if err != nil {
			return "", err
		}
		header := printableHeader("outputs", out.Name)
		text += header + tabtext
		if out.Name != "ostent" {
			ostent.AddExporter(header, printableConfigText(tabtext))
		}
	}

	return printableConfigText(text), nil
}

func printableHeader(a, b string) string { return fmt.Sprintf("    [%s.%s]\n", a, b) }

func printableTable(in1 interface{}, in2 *printInput, in3 *printFilter) (string, error) {
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
	return strings.Join(lines, "\n"), nil
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

func normalize(cf string, tab *ast.Table) error {
	var (
		ins  = &ast.Table{Fields: make(map[string]interface{})}
		outs = &ast.Table{Fields: make(map[string]interface{})}
	)

	if v, ok := tab.Fields["inputs"]; !ok {
		tab.Fields["inputs"] = ins
	} else if tfi, ok := v.(*ast.Table); !ok || tfi == nil || tfi.Fields == nil {
		tab.Fields["inputs"] = ins
	} else {
		ins = tfi
	}

	if v, ok := tab.Fields["outputs"]; !ok {
		tab.Fields["outputs"] = outs
	} else if tfo, ok := v.(*ast.Table); !ok || tfo == nil || tfo.Fields == nil {
		tab.Fields["outputs"] = outs
	} else {
		outs = tfo
	}

	for iname, ctext := range map[string]string{
		"cpu":             ``,
		"disk":            `ignore_fs = ["tmpfs", "devtmpfs"]`,
		"mem":             ``,
		"net_ostent":      ``,
		"procstat_ostent": ``,
		"swap":            ``,
		"system_ostent":   `interval = "1s"`,
	} {
		if _, ok := ins.Fields[iname]; ok {
			continue
		}
		ctab := &ast.Table{Fields: make(map[string]interface{})}
		if ctext != "" {
			var err error
			if ctab, err = config.ParseContents([]byte(ctext)); err != nil {
				return err
			}
		}
		ins.Fields[iname] = ctab
	}

	for oname, ctext := range map[string]string{
		"ostent": ``,
	} {
		if _, ok := outs.Fields[oname]; ok {
			continue
		}
		ctab := &ast.Table{Fields: make(map[string]interface{})}
		if ctext != "" {
			var err error
			if ctab, err = config.ParseContents([]byte(ctext)); err != nil {
				return err
			}
		}
		outs.Fields[oname] = ctab
	}

	deleteDisable(cf, "inputs", ins)
	deleteDisable(cf, "inputs", outs)

	var nonostentOutputs int
	for name := range outs.Fields {
		if name != "ostent" {
			nonostentOutputs++
		}
	}
	if nonostentOutputs > 0 {
		commondrop, err := config.ParseContents([]byte(`
namedrop = ["procstat", "procstat_ostent"]
[tagdrop]
    kind = ["system_ostent_runtime"]
`))
		if err != nil {
			return err
		}
		for name, value := range outs.Fields {
			if name != "ostent" {
				setfield(value, "namedrop", commondrop.Fields["namedrop"])
				setfield(value, "tagdrop", commondrop.Fields["tagdrop"])
			}
		}
	}
	return nil
}

func setfield(value interface{}, key string, set interface{}) {
	vtab, ok := value.(*ast.Table)
	if !ok {
		return
	}
	_, ok = vtab.Fields[key]
	if ok {
		return
	}
	vtab.Fields[key] = set
}

func deleteDisable(cf string, tname string, tab *ast.Table) {
	for name, value := range tab.Fields {
		vtab, ok := value.(*ast.Table)
		if !ok {
			continue
		}
		k := "disabled" // the proper key for deletion
		bv, ok := vtab.Fields[k]
		if !ok {
			k = "disable" // an alias key for deletion
			bv, ok = vtab.Fields[k]
		}
		if !ok {
			continue
		}
		kpath := fmt.Sprintf("%s.%s.%s", tname, name, k)

		// default config (cf == "") should not have `disabled` entries
		// so log calls with empty cf won't be caused.

		bkv, ok := bv.(*ast.KeyValue)
		if !ok {
			log.Printf("%s: Warn: %s value is of wrong type\n", cf, kpath)
			continue
		}
		bb, ok := bkv.Value.(*ast.Boolean)
		if !ok {
			log.Printf("%s: Warn: %s value is not a boolean\n", cf, kpath)
			continue
		}
		if b, err := bb.Boolean(); err == nil && b {
			delete(tab.Fields, name)
			log.Printf("%s: Info: [%s.%s] is disabled\n", cf, tname, name)
		}
	}
}

func readConfig(rconfig *config.Config) (*ast.Table, error) {
	var tab *ast.Table
	cf := viper.ConfigFileUsed()
	if cf != "" {
		// fmt.Printf("Using config file parsed:\n%#v\n", viper.AllSettings())

		text, err := ioutil.ReadFile(cf)
		if err != nil {
			return nil, err
		}
		tab, err = config.ParseContents(text)
		if err != nil {
			return nil, err
		}
	}
	if tab == nil {
		tab = &ast.Table{Fields: make(map[string]interface{})}
	} else if tab.Fields == nil {
		tab.Fields = make(map[string]interface{})
	}
	if err := normalize(cf, tab); err != nil {
		return nil, err
	}
	return tab, nil
}
