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
	_ "github.com/influxdata/telegraf/plugins/inputs/system"
	"github.com/influxdata/telegraf/plugins/outputs"
	_ "github.com/influxdata/telegraf/plugins/outputs/graphite"
	_ "github.com/influxdata/telegraf/plugins/outputs/influxdb"
	_ "github.com/influxdata/telegraf/plugins/outputs/librato"
	"github.com/influxdata/telegraf/plugins/serializers"

	"github.com/influxdata/toml"
	"github.com/influxdata/toml/ast"

	"github.com/ostrost/ostent/internal"
	internal_models "github.com/ostrost/ostent/internal/models"
	_ "github.com/ostrost/ostent/internal/plugins/outputs/ostent" // "ostent" output
	_ "github.com/ostrost/ostent/system_ostent"                   // "system_ostent" input
)

var config = struct {
	UnmarshalTable func(*ast.Table, interface{}) error
}{UnmarshalTable: toml.UnmarshalTable}

var (
	// envVarRe is a regex to find environment variables in the config file
	envVarRe = regexp.MustCompile(`\$\w+`)
)

// Config specifies the URL/user/password for the database that the agent
// will be logging to, as well as all the plugins that the user has
// specified
type Config struct {
	Tags map[string]string

	Agent   *AgentConfig
	Inputs  []*internal_models.RunningInput
	Outputs []*internal_models.RunningOutput
}

func NewConfig() *Config {
	c := &Config{
		Agent: &AgentConfig{
			Interval:      internal.Duration{Duration: 10 * time.Second},
			FlushInterval: internal.Duration{Duration: 10 * time.Second},
		},
	}
	return c
}

type AgentConfig struct {
	// Interval at which to gather information
	Interval internal.Duration

	// FlushInterval is the Interval at which to flush data
	FlushInterval internal.Duration
}

// LoadConfig loads the given config file and applies it to c
func (c *Config) LoadConfig(path string) error {
	var err error
	if path == "" {
		return fmt.Errorf("No config file specified")
		getDefaultConfigPath := func() (string, error) { return "", nil }
		if path, err = getDefaultConfigPath(); err != nil {
			return err
		}
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

	ro := internal_models.NewRunningOutput(name, output, outputConfig)
	c.Outputs = append(c.Outputs, ro)
	return nil
}

func (c *Config) addInput(name string, table *ast.Table) error {
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

	if err := f.CompileFilter(); err != nil {
		return f, err
	}

	delete(tbl.Fields, "namedrop")
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

	delete(tbl.Fields, "interval")
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
	/*
		// Outputs don't support FieldDrop/FieldPass, so set to NameDrop/NamePass
		if len(oc.Filter.FieldDrop) > 0 {
			oc.Filter.NameDrop = oc.Filter.FieldDrop
		}
		if len(oc.Filter.FieldPass) > 0 {
			oc.Filter.NamePass = oc.Filter.FieldPass
		}
	*/
	return oc, nil
}

func (c *Config) LoadInterface(path string, in interface{}) error {
	text, err := toml.Marshal(in)
	if err != nil {
		return err
	}
	log.Printf("#%s TOML formatted:\n%s", path, text)
	tbl, err := parseContents(text)
	if err != nil {
		return err
	}
	return c.LoadTable(path, tbl)
}
