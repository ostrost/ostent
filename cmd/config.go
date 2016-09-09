package cmd

import (
	"io/ioutil"

	"github.com/influxdata/toml/ast"
	"github.com/spf13/viper"

	"github.com/ostrost/ostent/internal/config"
)

// TODO intervals from flag(s)
func normalize(tab *ast.Table) error {
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
		"system_ostent":   ``,
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

	deleteDisable(ins)
	deleteDisable(outs)

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

func deleteDisable(tab *ast.Table) {
	for name, value := range tab.Fields {
		if vtab, ok := value.(*ast.Table); ok {
			if bv, ok := vtab.Fields["disable"]; ok {
				if bkv, ok := bv.(*ast.KeyValue); ok {
					if bb, ok := bkv.Value.(*ast.Boolean); ok {
						if b, err := bb.Boolean(); err == nil && b {
							delete(tab.Fields, name)
						}
					}
				}
			}
		}
	}
}

func readConfig(rconfig *config.Config) (*ast.Table, error) {
	var tab *ast.Table
	if cf := viper.ConfigFileUsed(); cf != "" {
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
	if err := normalize(tab); err != nil {
		return nil, err
	}
	return tab, nil
}
