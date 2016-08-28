package config

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/influxdata/telegraf/plugins/inputs"
	"github.com/influxdata/telegraf/plugins/outputs"
	"github.com/influxdata/telegraf/plugins/serializers"

	"github.com/influxdata/toml"
	"github.com/influxdata/toml/ast"

	"github.com/ostrost/ostent/internal"
	internal_models "github.com/ostrost/ostent/internal/models"
)

var config = struct {
	UnmarshalTable func(*ast.Table, interface{}) error
}{UnmarshalTable: toml.UnmarshalTable}

var (
	// Default input plugins
	inputDefaults = []string{"cpu", "mem", "swap", "system", "kernel",
		"processes", "disk", "diskio"}

	// Default output plugins
	outputDefaults = []string{"influxdb"}

	// envVarRe is a regex to find environment variables in the config file
	envVarRe = regexp.MustCompile(`\$\w+`)
)

// Config specifies the URL/user/password for the database that telegraf
// will be logging to, as well as all the plugins that the user has
// specified
type Config struct {
	Tags          map[string]string
	InputFilters  []string
	OutputFilters []string

	Agent   *AgentConfig
	Inputs  []*internal_models.RunningInput
	Outputs []*internal_models.RunningOutput
}

func NewConfig() *Config {
	c := &Config{
		// Agent defaults:
		Agent: &AgentConfig{
			Interval:      internal.Duration{Duration: 10 * time.Second},
			FlushInterval: internal.Duration{Duration: 10 * time.Second},
		},

		Tags:          make(map[string]string),
		InputFilters:  make([]string, 0),
		OutputFilters: make([]string, 0),
	}
	return c
}

type AgentConfig struct {
	// Interval at which to gather information
	Interval internal.Duration

	// By default, precision will be set to the same timestamp order as the
	// collection interval, with the maximum being 1s.
	//   ie, when interval = "10s", precision will be "1s"
	//       when interval = "250ms", precision will be "1ms"
	// Precision will NOT be used for service inputs. It is up to each individual
	// service input to set the timestamp at the appropriate precision.
	Precision internal.Duration

	// FlushInterval is the Interval at which to flush data
	FlushInterval internal.Duration

	// MetricBatchSize is the maximum number of metrics that is wrote to an
	// output plugin in one call.
	MetricBatchSize int

	// MetricBufferLimit is the max number of metrics that each output plugin
	// will cache. The buffer is cleared when a successful write occurs. When
	// full, the oldest metrics will be overwritten. This number should be a
	// multiple of MetricBatchSize. Due to current implementation, this could
	// not be less than 2 times MetricBatchSize.
	MetricBufferLimit int

	// Debug is the option for running in debug mode
	Debug bool

	// Quiet is the option for running in quiet mode
	Quiet        bool
	Hostname     string
	OmitHostname bool
}

func sliceContains(name string, list []string) bool {
	for _, b := range list {
		if b == name {
			return true
		}
	}
	return false
}

// LoadConfig loads the given config file and applies it to c
func (c *Config) LoadConfig(path string) error {
	var err error
	if path == "" {
		return fmt.Errorf("No config file specified")
		/*
			if path, err = getDefaultConfigPath(); err != nil {
				return err
			}
		*/
	}
	tbl, err := parseFile(path)
	if err != nil {
		return fmt.Errorf("Error parsing %s, %s", path, err)
	}

	return c.LoadTable(path, tbl)
}

func (c *Config) LoadTable(path string, tbl *ast.Table) error {
	var err error

	// Parse tags tables first:
	for _, tableName := range []string{"tags", "global_tags"} {
		if val, ok := tbl.Fields[tableName]; ok {
			subTable, ok := val.(*ast.Table)
			if !ok {
				return fmt.Errorf("%s: invalid configuration", path)
			}
			if err = config.UnmarshalTable(subTable, c.Tags); err != nil {
				log.Printf("Could not parse [global_tags] config\n")
				return fmt.Errorf("Error parsing %s, %s", path, err)
			}
		}
	}

	// Parse agent table:
	if val, ok := tbl.Fields["agent"]; ok {
		subTable, ok := val.(*ast.Table)
		if !ok {
			return fmt.Errorf("%s: invalid configuration", path)
		}
		if err = config.UnmarshalTable(subTable, c.Agent); err != nil {
			log.Printf("Could not parse [agent] config\n")
			return fmt.Errorf("Error parsing %s, %s", path, err)
		}
	}

	// Parse all the rest of the plugins:
	for name, val := range tbl.Fields {
		subTable, ok := val.(*ast.Table)
		if !ok {
			return fmt.Errorf("%s: invalid configuration", path)
		}

		switch name {
		case "agent", "global_tags", "tags":
		case "outputs":
			for pluginName, pluginVal := range subTable.Fields {
				switch pluginSubTable := pluginVal.(type) {
				case *ast.Table:
					if err = c.addOutput(pluginName, pluginSubTable); err != nil {
						return fmt.Errorf("Error parsing %s, %s", path, err)
					}
				case []*ast.Table:
					for _, t := range pluginSubTable {
						if err = c.addOutput(pluginName, t); err != nil {
							return fmt.Errorf("Error parsing %s, %s", path, err)
						}
					}
				default:
					return fmt.Errorf("Unsupported config format: %s, file %s",
						pluginName, path)
				}
			}
		case "inputs", "plugins":
			for pluginName, pluginVal := range subTable.Fields {
				switch pluginSubTable := pluginVal.(type) {
				case *ast.Table:
					if err = c.addInput(pluginName, pluginSubTable); err != nil {
						return fmt.Errorf("Error parsing %s, %s", path, err)
					}
				case []*ast.Table:
					for _, t := range pluginSubTable {
						if err = c.addInput(pluginName, t); err != nil {
							return fmt.Errorf("Error parsing %s, %s", path, err)
						}
					}
				default:
					return fmt.Errorf("Unsupported config format: %s, file %s",
						pluginName, path)
				}
			}
			// Assume it's an input input for legacy config file support if no other
			// identifiers are present
		default:
			if err = c.addInput(name, subTable); err != nil {
				return fmt.Errorf("Error parsing %s, %s", path, err)
			}
		}
	}
	return nil
}

// trimBOM trims the Byte-Order-Marks from the beginning of the file.
// this is for Windows compatability only.
// see https://github.com/influxdata/telegraf/issues/1378
func trimBOM(f []byte) []byte {
	return bytes.TrimPrefix(f, []byte("\xef\xbb\xbf"))
}

// parseFile loads a TOML configuration from a provided path and
// returns the AST produced from the TOML parser. When loading the file, it
// will find environment variables and replace them.
func parseFile(fpath string) (*ast.Table, error) {
	contents, err := ioutil.ReadFile(fpath)
	if err != nil {
		return nil, err
	}
	return parseContents(contents)
}

func parseContents(contents []byte) (*ast.Table, error) {
	// ugh windows why
	contents = trimBOM(contents)

	env_vars := envVarRe.FindAll(contents, -1)
	for _, env_var := range env_vars {
		env_val := os.Getenv(strings.TrimPrefix(string(env_var), "$"))
		if env_val != "" {
			contents = bytes.Replace(contents, env_var, []byte(env_val), 1)
		}
	}

	return toml.Parse(contents)
}

func (c *Config) addOutput(name string, table *ast.Table) error {
	if len(c.OutputFilters) > 0 && !sliceContains(name, c.OutputFilters) {
		return nil
	}
	creator, ok := outputs.Outputs[name]
	if !ok {
		return fmt.Errorf("Undefined but requested output: %s", name)
	}
	output := creator()

	// If the output has a SetSerializer function, then this means it can write
	// arbitrary types of output, so build the serializer and set it.
	switch t := output.(type) {
	case serializers.SerializerOutput:
		serializer, err := buildSerializer(name, table)
		if err != nil {
			return err
		}
		if serializer == nil {
			return fmt.Errorf("Serializer is nil")
		}
		t.SetSerializer(serializer)
	}

	outputConfig, err := buildOutput(name, table)
	if err != nil {
		return err
	}

	if err := config.UnmarshalTable(table, output); err != nil {
		return err
	}

	ro := internal_models.NewRunningOutput(name, output, outputConfig,
		c.Agent.MetricBatchSize, c.Agent.MetricBufferLimit)
	c.Outputs = append(c.Outputs, ro)
	return nil
}

func (c *Config) addInput(name string, table *ast.Table) error {
	if len(c.InputFilters) > 0 && !sliceContains(name, c.InputFilters) {
		return nil
	}
	// Legacy support renaming io input to diskio
	if name == "io" {
		name = "diskio"
	}

	creator, ok := inputs.Inputs[name]
	if !ok {
		return fmt.Errorf("Undefined but requested input: %s", name)
	}
	input := creator()
	pluginConfig, err := buildInput(name, table)
	if err != nil {
		return err
	}

	if err := config.UnmarshalTable(table, input); err != nil {
		return err
	}

	rp := &internal_models.RunningInput{
		Name:   name,
		Input:  input,
		Config: pluginConfig,
	}
	c.Inputs = append(c.Inputs, rp)
	return nil
}

// buildFilter builds a Filter
// (tagpass/tagdrop/namepass/namedrop/fieldpass/fielddrop) to
// be inserted into the internal_models.OutputConfig/internal_models.InputConfig
// to be used for glob filtering on tags and measurements
func buildFilter(tbl *ast.Table) (internal_models.Filter, error) {
	f := internal_models.Filter{}

	if node, ok := tbl.Fields["namepass"]; ok {
		if kv, ok := node.(*ast.KeyValue); ok {
			if ary, ok := kv.Value.(*ast.Array); ok {
				for _, elem := range ary.Value {
					if str, ok := elem.(*ast.String); ok {
						f.NamePass = append(f.NamePass, str.Value)
						f.IsActive = true
					}
				}
			}
		}
	}

	if node, ok := tbl.Fields["namedrop"]; ok {
		if kv, ok := node.(*ast.KeyValue); ok {
			if ary, ok := kv.Value.(*ast.Array); ok {
				for _, elem := range ary.Value {
					if str, ok := elem.(*ast.String); ok {
						f.NameDrop = append(f.NameDrop, str.Value)
						f.IsActive = true
					}
				}
			}
		}
	}

	fields := []string{"pass", "fieldpass"}
	for _, field := range fields {
		if node, ok := tbl.Fields[field]; ok {
			if kv, ok := node.(*ast.KeyValue); ok {
				if ary, ok := kv.Value.(*ast.Array); ok {
					for _, elem := range ary.Value {
						if str, ok := elem.(*ast.String); ok {
							f.FieldPass = append(f.FieldPass, str.Value)
							f.IsActive = true
						}
					}
				}
			}
		}
	}

	fields = []string{"drop", "fielddrop"}
	for _, field := range fields {
		if node, ok := tbl.Fields[field]; ok {
			if kv, ok := node.(*ast.KeyValue); ok {
				if ary, ok := kv.Value.(*ast.Array); ok {
					for _, elem := range ary.Value {
						if str, ok := elem.(*ast.String); ok {
							f.FieldDrop = append(f.FieldDrop, str.Value)
							f.IsActive = true
						}
					}
				}
			}
		}
	}

	if node, ok := tbl.Fields["tagpass"]; ok {
		if subtbl, ok := node.(*ast.Table); ok {
			for name, val := range subtbl.Fields {
				if kv, ok := val.(*ast.KeyValue); ok {
					tagfilter := &internal_models.TagFilter{Name: name}
					if ary, ok := kv.Value.(*ast.Array); ok {
						for _, elem := range ary.Value {
							if str, ok := elem.(*ast.String); ok {
								tagfilter.Filter = append(tagfilter.Filter, str.Value)
							}
						}
					}
					f.TagPass = append(f.TagPass, *tagfilter)
					f.IsActive = true
				}
			}
		}
	}

	if node, ok := tbl.Fields["tagdrop"]; ok {
		if subtbl, ok := node.(*ast.Table); ok {
			for name, val := range subtbl.Fields {
				if kv, ok := val.(*ast.KeyValue); ok {
					tagfilter := &internal_models.TagFilter{Name: name}
					if ary, ok := kv.Value.(*ast.Array); ok {
						for _, elem := range ary.Value {
							if str, ok := elem.(*ast.String); ok {
								tagfilter.Filter = append(tagfilter.Filter, str.Value)
							}
						}
					}
					f.TagDrop = append(f.TagDrop, *tagfilter)
					f.IsActive = true
				}
			}
		}
	}

	if node, ok := tbl.Fields["tagexclude"]; ok {
		if kv, ok := node.(*ast.KeyValue); ok {
			if ary, ok := kv.Value.(*ast.Array); ok {
				for _, elem := range ary.Value {
					if str, ok := elem.(*ast.String); ok {
						f.TagExclude = append(f.TagExclude, str.Value)
					}
				}
			}
		}
	}

	if node, ok := tbl.Fields["taginclude"]; ok {
		if kv, ok := node.(*ast.KeyValue); ok {
			if ary, ok := kv.Value.(*ast.Array); ok {
				for _, elem := range ary.Value {
					if str, ok := elem.(*ast.String); ok {
						f.TagInclude = append(f.TagInclude, str.Value)
					}
				}
			}
		}
	}
	if err := f.CompileFilter(); err != nil {
		return f, err
	}

	delete(tbl.Fields, "namedrop")
	delete(tbl.Fields, "namepass")
	delete(tbl.Fields, "fielddrop")
	delete(tbl.Fields, "fieldpass")
	delete(tbl.Fields, "drop")
	delete(tbl.Fields, "pass")
	delete(tbl.Fields, "tagdrop")
	delete(tbl.Fields, "tagpass")
	delete(tbl.Fields, "tagexclude")
	delete(tbl.Fields, "taginclude")
	return f, nil
}

// buildSerializer grabs the necessary entries from the ast.Table for creating
// a serializers.Serializer object, and creates it, which can then be added onto
// an Output object.
func buildSerializer(name string, tbl *ast.Table) (serializers.Serializer, error) {
	return serializers.NewSerializer(&serializers.Config{
		DataFormat: "graphite",
	})
}

func buildInput(name string, tbl *ast.Table) (*internal_models.InputConfig, error) {
	cp := &internal_models.InputConfig{Name: name}
	if node, ok := tbl.Fields["interval"]; ok {
		if kv, ok := node.(*ast.KeyValue); ok {
			if str, ok := kv.Value.(*ast.String); ok {
				dur, err := time.ParseDuration(str.Value)
				if err != nil {
					return nil, err
				}

				cp.Interval = dur
			}
		}
	}

	if node, ok := tbl.Fields["name_prefix"]; ok {
		if kv, ok := node.(*ast.KeyValue); ok {
			if str, ok := kv.Value.(*ast.String); ok {
				cp.MeasurementPrefix = str.Value
			}
		}
	}

	if node, ok := tbl.Fields["name_suffix"]; ok {
		if kv, ok := node.(*ast.KeyValue); ok {
			if str, ok := kv.Value.(*ast.String); ok {
				cp.MeasurementSuffix = str.Value
			}
		}
	}

	if node, ok := tbl.Fields["name_override"]; ok {
		if kv, ok := node.(*ast.KeyValue); ok {
			if str, ok := kv.Value.(*ast.String); ok {
				cp.NameOverride = str.Value
			}
		}
	}

	cp.Tags = make(map[string]string)
	if node, ok := tbl.Fields["tags"]; ok {
		if subtbl, ok := node.(*ast.Table); ok {
			if err := config.UnmarshalTable(subtbl, cp.Tags); err != nil {
				log.Printf("Could not parse tags for input %s\n", name)
			}
		}
	}

	delete(tbl.Fields, "name_prefix")
	delete(tbl.Fields, "name_suffix")
	delete(tbl.Fields, "name_override")
	delete(tbl.Fields, "interval")
	delete(tbl.Fields, "tags")
	var err error
	cp.Filter, err = buildFilter(tbl)
	if err != nil {
		return cp, err
	}
	return cp, nil
}

// buildOutput parses output specific items from the ast.Table,
// builds the filter and returns an
// internal_models.OutputConfig to be inserted into internal_models.RunningInput
// Note: error exists in the return for future calls that might require error
func buildOutput(name string, tbl *ast.Table) (*internal_models.OutputConfig, error) {
	filter, err := buildFilter(tbl)
	if err != nil {
		return nil, err
	}
	oc := &internal_models.OutputConfig{
		Name:   name,
		Filter: filter,
	}
	// Outputs don't support FieldDrop/FieldPass, so set to NameDrop/NamePass
	if len(oc.Filter.FieldDrop) > 0 {
		oc.Filter.NameDrop = oc.Filter.FieldDrop
	}
	if len(oc.Filter.FieldPass) > 0 {
		oc.Filter.NamePass = oc.Filter.FieldPass
	}
	return oc, nil
}

func (c *Config) LoadInterface(path string, in interface{}) error {
	text, err := toml.Marshal(in)
	if err != nil {
		return err
	}
	lines := strings.Split(string(text), "\n")
	for _, replace := range [][2]string{
		{"password=", "********"},
		{"api_token=", "****************"},
	} {
		for i := range lines {
			if strings.HasPrefix(lines[i], replace[0]) {
				lines[i] = fmt.Sprintf("%s=\"%s\"", replace[0], replace[1])
			}
		}
	}
	log.Printf("#%s.toml:\n%s", path, strings.Join(lines, "\n"))
	tbl, err := parseContents(text)
	if err != nil {
		return err
	}
	return c.LoadTable(path, tbl)
}
