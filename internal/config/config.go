package config

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/influxdata/telegraf/plugins/inputs"
	"github.com/influxdata/telegraf/plugins/outputs"
	"github.com/influxdata/telegraf/plugins/serializers"

	"github.com/influxdata/toml"
	"github.com/influxdata/toml/ast"

	"github.com/ostrost/ostent/internal"
	"github.com/ostrost/ostent/internal/buffer"
	internal_models "github.com/ostrost/ostent/internal/models"
)

var (
	// envVarRe is a regex to find environment variables in the config file
	envVarRe = regexp.MustCompile(`\$\w+`)
)

// Config specifies the URL/user/password for the database that the agent
// will be logging to, as well as all the plugins that the user has
// specified
type Config struct {
	Agent   *AgentConfig
	Inputs  []*internal_models.RunningInput
	Outputs []*internal_models.RunningOutput
}

func NewConfig() *Config {
	return &Config{
		Agent: &AgentConfig{
			// values are defaults
			Interval:      internal.Duration{Duration: time.Second * 10},
			FlushInterval: internal.Duration{Duration: time.Second * 10},
		},
	}
}

type AgentConfig struct {
	Interval      internal.Duration
	FlushInterval internal.Duration
}

// trimBOM trims the Byte-Order-Marks from the beginning of the file.
// this is for Windows compatability only.
// see https://github.com/influxdata/telegraf/issues/1378
func trimBOM(f []byte) []byte {
	return bytes.TrimPrefix(f, []byte("\xef\xbb\xbf"))
}

func parse(contents []byte) (*ast.Table, error) {
	contents = trimBOM(contents)
	for _, dword := range envVarRe.FindAll(contents, -1) {
		if val := os.Getenv(string(dword[1:])); val != "" {
			contents = bytes.Replace(contents, dword, []byte(val), 1)
		}
	}
	return toml.Parse(contents)
}

func (c *Config) LoadConfig() error {
	tbl, err := parse([]byte(`
[agent]
  interval = "1s"
  flushInterval = "1s"
[[inputs.system_ostent]]
  interval = "1s"
# [[outputs.file]]
[[outputs.ostent]]
`))
	if err != nil {
		return err
	}
	if val, ok := tbl.Fields["agent"]; ok {
		subt, ok := val.(*ast.Table)
		if !ok {
			return fmt.Errorf("Cannot parse config")
		}
		if err := toml.UnmarshalTable(subt, c.Agent); err != nil {
			return fmt.Errorf("Cannot parse config: [agent] section: %s", err)
		}
	}

	for name, val := range tbl.Fields {
		subt, ok := val.(*ast.Table)
		if !ok {
			return fmt.Errorf("Cannot parse config")
		}
		if name == "outputs" {
			for pname, pval := range subt.Fields {
				switch psubt := pval.(type) {
				case *ast.Table:
					if err := c.AddOutput(pname, psubt); err != nil {
						return fmt.Errorf("Parse error: %s", err)
					}
				case []*ast.Table:
					for _, t := range psubt {
						if err := c.AddOutput(pname, t); err != nil {
							return fmt.Errorf("Parse error: %s", err)
						}
					}
				default:
					return fmt.Errorf("Unsupported type in config: [%s] section", pname)
				}
			}
		} else if name == "inputs" {
			for pname, pval := range subt.Fields {
				switch psubt := pval.(type) {
				case *ast.Table:
					if err := c.AddInput(pname, psubt); err != nil {
						return fmt.Errorf("Parse error: %s", err)
					}
				case []*ast.Table:
					for _, t := range psubt {
						if err := c.AddInput(pname, t); err != nil {
							return fmt.Errorf("Parse error: %s", err)
						}
					}
				default:
					return fmt.Errorf("Unsupported type in config: [%s] section", pname)
				}
			}
		}
	}
	return err
}

func makeSerializer(name string) (serializers.Serializer, error) {
	return serializers.NewSerializer(&serializers.Config{
		DataFormat: "graphite",
	})
}

func (c *Config) AddOutput(name string, table *ast.Table) error {
	create, ok := outputs.Outputs[name]
	if !ok {
		return fmt.Errorf("Unknown output by name %q", name)
	}
	out := create()

	if ser, ok := out.(serializers.SerializerOutput); ok {
		newSer, err := makeSerializer(name)
		if err != nil {
			return err
		}
		if newSer == nil {
			return fmt.Errorf("Serializer is nil")
		}
		ser.SetSerializer(newSer)
	}

	bbs := 1000
	c.Outputs = append(c.Outputs, &internal_models.RunningOutput{
		Output:       out,
		Name:         name,
		BufBatchSize: bbs,
		Buf:          buffer.NewBuffer(bbs),
	})
	return nil
}

func (c *Config) AddInput(name string, table *ast.Table) error {
	create, ok := inputs.Inputs[name]
	if !ok {
		return fmt.Errorf("Unknown input by name %q", name)
	}
	in := create()
	pc, err := buildInput(name, table)
	if err != nil {
		return err
	}
	c.Inputs = append(c.Inputs, &internal_models.RunningInput{
		Input:  in,
		Name:   name,
		Config: pc,
	})
	return nil
}

func buildInput(name string, tbl *ast.Table) (*internal_models.InputConfig, error) {
	conf := &internal_models.InputConfig{Name: name, Interval: time.Second * 10}
	if node, ok := tbl.Fields["interval"]; ok {
		if kv, ok := node.(*ast.KeyValue); ok {
			if sv, ok := kv.Value.(*ast.String); ok {
				d, err := time.ParseDuration(sv.Value)
				if err != nil {
					return nil, err
				}
				conf.Interval = d
			}
		}
	}
	return conf, nil
}
