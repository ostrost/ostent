package templatepipe

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"sort"
	"strings"
	templatetext "text/template"

	// "golang.org/x/net/html"
)

func Main(htmlFuncs, jsxlFuncs map[string]interface{}) {
	var (
		outputFile       string
		htmltemplateFile string
	)
	flag.StringVar(&outputFile, "o", "", "Output file")
	flag.StringVar(&outputFile, "output", "", "Output file")
	flag.StringVar(&htmltemplateFile, "html", "", "The html template file to parse")
	flag.StringVar(&htmltemplateFile, "htmltemplate", "", "The html template file to parse")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage:\n  %s [options] [filename.ace]\n\nOptions:\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(64)
	}
	flag.Parse()

	inputFile := flag.Arg(0)
	if inputFile == "" {
		fmt.Fprintf(os.Stderr, "No template specified.\n")
		flag.Usage() // exits 64
		return
	}

	check := func(err error) {
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	jstemplate, err := templatetext.ParseFiles(inputFile)
	check(err)

	// defines, err := templatetext.ParseFiles(htmltemplateFile)
	defines, err := templatetext.New(htmltemplateFile).Funcs(htmlFuncs).
		ParseFiles(htmltemplateFile)
	check(err)

	definesonly := templatetext.New("jsdefines").
		Funcs(templatetext.FuncMap(jsxlFuncs))

	definesTemplates := SortableTemplates(defines.Templates())
	sort.Stable(definesTemplates)

	for _, t := range definesTemplates {
		_, err := definesonly.AddParseTree(t.Name(), t.Tree)
		check(err)
	}

	jdata := struct{ Defines []Define }{}
	for _, t := range definesTemplates {
		names := strings.Split(t.Name(), "::")
		tname := names[len(names)-1]
		if strings.HasPrefix(tname, "define_") {
			define, err := MakeDefine(definesonly, tname, t.Name())
			check(err)
			jdata.Defines = append(jdata.Defines, define)
		}
	}
	buf := new(bytes.Buffer)
	check(jstemplate.Execute(buf, jdata))
	check(WriteFile(outputFile, buf.String()))
}

func JSX2HTML(buf *bytes.Buffer) (string, error) {
	i, keys := 0, make([]string, len(JSXAttributeRewrites))
	for k := range JSXAttributeRewrites {
		keys[i] = k
		i++
	}
	sort.Strings(keys)

	s := buf.String()
	for _, k := range keys {
		s = strings.Replace(s, k, JSXAttributeRewrites[k], -1)
	}
	return s, nil
	/*
		node, err := html.Parse(buf)
		if err != nil {
			return "", err
		}
		JSXAttributes(node)
		buf := new(bytes.Buffer)
		html.Render(buf, node)
		return buf.String(), nil
	// */
}

// JSXAttributeRewrites is a map to jsx-compat attibute names.
var JSXAttributeRewrites = map[string]string{
	"colspan":   "colSpan",
	"class":     "lcassName",
	"lcassName": "className",
}

/*
// JSXAttributes replaces node and it's children attributes with rewrites from JSXAttributeRewrites.
func JSXAttributes(node *html.Node) {
	for i := range node.Attr {
		if nv, ok := JSXAttributeRewrites[node.Attr[i].Key]; ok {
			node.Attr[i].Key = nv
		}
	}
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		JSXAttributes(c)
	}
} // */

type Define struct {
	ShortName  string
	Iterable   string
	NeedList   bool
	UsesParams bool
	JSX        string
}

// SortableTemplates is for just sorting.
type SortableTemplates []*templatetext.Template

func (st SortableTemplates) Len() int           { return len(st) }
func (st SortableTemplates) Less(i, j int) bool { return st[i].Name() < st[j].Name() }
func (st SortableTemplates) Swap(i, j int)      { st[i], st[j] = st[j], st[i] }

func MakeDefine(definesonly *templatetext.Template, shortname, fullname string) (Define, error) {
	define := Define{ShortName: shortname}
	t, err := definesonly.Clone()
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
	define.JSX, err = JSX2HTML(buf)
	return define, err
}

var vtype = reflect.TypeOf(Nota(nil))

func CurlyNotMethod(parent, key, full string) interface{} {
	if _, ok := vtype.MethodByName(key); ok {
		return nil
	}
	return CurlyX(parent, key, full)
}

// WriteFile is ioutil.WriteFile if filename is not "",
// otherwise it's as if filename was /dev/stdout.
func WriteFile(filename, data string) error {
	bytedata := []byte(data)
	if filename != "" {
		return ioutil.WriteFile(filename, bytedata, 0644)
	}
	_, err := os.Stdout.Write(bytedata)
	return err
}
