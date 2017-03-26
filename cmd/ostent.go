package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"unicode"

	"github.com/influxdata/toml"
	"github.com/influxdata/toml/ast"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

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

var fv = &flagValues{
	FlagSet:  RootCmd.Flags(),
	bindPort: portFlag(8050), // default 8050
	interval: intervalFlag(config.NewConfig().Agent.Interval),
}

func initFlags() {
	persistentPreRuns.add(versionRun)
	fv.BoolVar(&versionFlag, "version", false, "print version and exit")

	fv.StringVar(&fv.bind, "bind", "",
		fmt.Sprintf("server bind address (default %q)", fv.bind))
	fv.Var(&fv.bindPort, "bind-port",
		fmt.Sprintf("server bind port (default %d)", fv.bindPort))
	fv.Var(&fv.interval, "interval",
		fmt.Sprintf("metrics collection interval (default %s)", fv.interval))
}

type (
	portFlag     int
	intervalFlag internal.Duration
	flagValues   struct {
		*pflag.FlagSet

		bind     string
		bindPort portFlag
		interval intervalFlag
	}
)

// String is of fmt.Stringer interface.
func (iv intervalFlag) String() string { return iv.Duration.String() }
func (iv intervalFlag) Type() string   { return "duration" }
func (iv *intervalFlag) Set(input string) error {
	var in internal.Duration
	q := []byte(fmt.Sprintf("%q", input))
	if err := in.UnmarshalTOML(q); err != nil {
		return err
	}
	*iv = intervalFlag(in)
	return nil
}

// String is of fmt.Stringer interface.
func (pf portFlag) String() string { return strconv.Itoa(int(pf)) }
func (pf portFlag) Type() string   { return "int" }
func (pf *portFlag) Set(input string) error {
	p, err := net.LookupPort("tcp", input)
	if err != nil {
		return err
	}
	*pf = portFlag(p)
	return nil
}

type setonce struct {
	mutex sync.Mutex
	set   bool
	host  string
	port  int
}

// change compares host and port with previously passed values.
func (o *setonce) change(host string, port int) (string, int, bool) {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	if o.set {
		return o.host, o.port, o.host != host || o.port != port
	}
	o.host, o.port, o.set = host, port, true
	return host, port, false
}

// MainAgent is the main agent job.
func MainAgent(send chan chan string) {
	if err := mainAgent(send, new(setonce)); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func mainAgent(send <-chan chan string, fsthp *setonce) error {
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
		err := loadConfig(c, send, fsthp)
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

// nolint: gocyclo
func loadConfig(c *config.Config, send <-chan chan string, fsthp *setonce) error {
	c.Agent.Quiet = true // patch work

	tab, cf, err := readConfig(c)
	if err != nil {
		return err
	} else if tab != nil {
		if err2 := c.LoadTable("/runtime/config", tab); err2 != nil {
			return err2
		}
	}

	// fv is a global
	if f := fv.Lookup("interval"); f != nil && f.Changed {
		c.Agent.Interval = internal.Duration(fv.interval)
	}
	host := fv.bind
	if f := fv.Lookup("bind"); f == nil || !f.Changed {
		host = c.Agent.Bind
	}
	port := int(fv.bindPort)
	if f := fv.Lookup("bind-port"); (f == nil || !f.Changed) && c.Agent.BindPort != 0 {
		port = c.Agent.BindPort
	}

	var change bool
	newhost, newport := host, port
	host, port, change = fsthp.change(host, port)
	hp := net.JoinHostPort(host, strconv.Itoa(port))
	if change {
		newhp := net.JoinHostPort(newhost, strconv.Itoa(newport))
		log.Printf("%s: Warn: new bind address:port require process restart\n", cf)
		log.Printf("%s: Warn: new bind address:port: %q\n", cf, newhp)
		log.Printf("%s: Warn: old bind address:port: %q\n", cf, hp)
		log.Printf("%s: Warn: old bind address:port in effect until a restart\n", cf)
	}

	c.Agent.Bind, c.Agent.BindPort = host, port

	if receive, ok := <-send; ok && receive != nil {
		receive <- hp
	}

	text, err := printableConfig(c)
	if err != nil {
		return err
	}
	log.Printf("Effective runtime config:\n%s", text)
	return nil
}

func hasEsuffix(s, suffix string) bool {
	return strings.HasSuffix(s, "="+suffix) || strings.HasSuffix(s, " = "+suffix)
}
func hasprefixE(s, prefix string) bool {
	return strings.HasPrefix(s, prefix+"=") || strings.HasPrefix(s, prefix+" = ")
}

func printableConfigText(text string) string {
	lines := strings.Split(text, "\n")
	var newlines []string
rangelines:
	for i, line := range lines {

		for _, suffix := range []string{
			"0",
			`""`,
			"[]",
			`"0s"`,
			"false",
		} {
			if hasEsuffix(line, suffix) {
				continue rangelines
			}
		}
		for _, x := range []struct{ name, defvalue string }{
			{"bind_port", "8050"},
			{"round_interval", "true"},
			{"quiet", "true"},
		} {
			if hasprefixE(strings.TrimLeftFunc(line, unicode.IsSpace), x.name) &&
				(hasEsuffix(line, x.defvalue) /* || hasEsuffix(line, fmt.Sprintf("%q", x.defvalue)) */) {
				continue rangelines
			}
		}

		for _, replace := range [][2]string{
			{"password", `"PASSWORD"`},
			{"api_token", `"API_TOKEN"`},
		} {
			if j := strings.Index(line, replace[0]+"="); j != -1 {
				lines[i] = line[:j] + replace[0] + "=" + replace[1]
			} else if j := strings.Index(line, replace[0]+" = "); j != -1 {
				lines[i] = line[:j] + replace[0] + " = " + replace[1]
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
		text += printableHeader("inputs", in.Config.Name) + tabtext
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
	if !f.IsActive() {
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

// nolint: gocyclo
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

	const ostentOutput = "ostent"
	for oname, ctext := range map[string]string{
		ostentOutput: ``,
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
		if name != ostentOutput {
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
			if name != ostentOutput {
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

func readConfig(rconfig *config.Config) (*ast.Table, string, error) {
	var tab *ast.Table
	cf := viper.ConfigFileUsed()
	if cf != "" {
		log.Printf("%s config file to use\n", cf)

		text, err := ioutil.ReadFile(cf)
		if err != nil {
			return nil, cf, err
		}
		tab, err = config.ParseContents(text)
		if err != nil {
			return nil, cf, err
		}
	}
	if tab == nil {
		tab = &ast.Table{Fields: make(map[string]interface{})}
	} else if tab.Fields == nil {
		tab.Fields = make(map[string]interface{})
	}
	if err := normalize(cf, tab); err != nil {
		return nil, cf, err
	}
	return tab, cf, nil
}
