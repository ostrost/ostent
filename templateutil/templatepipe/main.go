package templatepipe

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"sort"
	"strings"
	"text/template"
)

func Convert(inputTemplateFile, definesFromFile string,
	htmlFuncs /*, jsxlFuncs */ map[string]interface{}, outputFile string) error {
	input, err := template.ParseFiles(inputTemplateFile)
	if err != nil {
		return err
	}

	definesFrom := template.New(definesFromFile).Delims("[[", "]]").Funcs(htmlFuncs)
	if definesFromFile == "-" {
		text, err2 := ioutil.ReadAll(os.Stdin)
		if err2 != nil {
			return err2
		}
		definesFrom, err = definesFrom.Parse(string(text))
	} else {
		definesFrom, err = definesFrom.ParseFiles(definesFromFile)
	}
	if err != nil {
		return err
	}

	// definesOnly will have just "define_" templates added in the tree.
	definesOnly := template.New("jsxdefines") // .Funcs(template.FuncMap(jsxlFuncs))

	definesTemplates := SortableTemplates(definesFrom.Templates())
	sort.Stable(definesTemplates)

	for _, t := range definesTemplates {
		if _, err := definesOnly.AddParseTree(t.Name(), t.Tree); err != nil {
			return err
		}
	}

	jdata := struct{ Defines []Define }{}
	for _, t := range definesTemplates {
		if tname := t.Name(); strings.HasPrefix(tname, "define_") {
			define, err := MakeDefine(definesOnly, tname, tname)
			if err != nil {
				return err
			}
			jdata.Defines = append(jdata.Defines, define)
		}
	}
	buf := new(bytes.Buffer)
	if err := input.Execute(buf, jdata); err != nil {
		return err
	}
	if outputFile == "" || outputFile == "-" {
		_, err := os.Stdout.Write(buf.Bytes())
		return err
	}
	return ioutil.WriteFile(outputFile, buf.Bytes(), 0644)
}

type Define struct {
	ShortName  string
	Iterable   string
	NeedList   bool
	UsesParams bool
	JSX        string
}

// SortableTemplates is for just sorting.
type SortableTemplates []*template.Template

func (st SortableTemplates) Len() int           { return len(st) }
func (st SortableTemplates) Less(i, j int) bool { return st[i].Name() < st[j].Name() }
func (st SortableTemplates) Swap(i, j int)      { st[i], st[j] = st[j], st[i] }

func MakeDefine(definesOnly *template.Template, shortname, fullname string) (Define, error) {
	define := Define{ShortName: shortname}
	t, err := definesOnly.Clone()
	if err != nil {
		return define, err
	}
	if t, err = t.Parse(fmt.Sprintf(`{{template %q .}}`, fullname)); err != nil {
		return define, err
	}

	data := Data(CurlyNotMethod, t)
	if nota, ok := data.(Nota); ok {
		for k, v := range nota["Data"].(Nota) {
			if k == "params" {
				define.UsesParams = true
			} else if k != "." {
				if define.Iterable != "" {
					return define, fmt.Errorf("Key %q is second: iterable already by %q", k, define.Iterable)
				}
				define.Iterable = k
				if n, ok := v.(Nota); ok {
					if _, ok := n["List"]; ok {
						define.NeedList = true
					}
				}
			}
		}
	}
	buf := new(bytes.Buffer)
	if err := t.Execute(buf, data); err != nil {
		return define, err
	}
	define.JSX = strings.Replace(buf.String(), "class=", "className=", -1)
	return define, nil
}

var vtype = reflect.TypeOf(Nota(nil))

func CurlyNotMethod(parent, key, full string) interface{} {
	if _, ok := vtype.MethodByName(key); ok {
		return nil
	}
	return CurlyX(parent, key, full)
}
