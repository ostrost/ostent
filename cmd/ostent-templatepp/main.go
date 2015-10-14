package main

import (
	"bytes"
	"flag"
	"fmt"
	templatehtml "html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	templatetext "text/template"
	"text/template/parse"

	"github.com/yosssi/ace"
	"golang.org/x/net/html"

	"github.com/ostrost/ostent/templateutil/templatefunc"
	"github.com/ostrost/ostent/templateutil/templatepipe"
)

func main() {
	var (
		outputFile    string
		definesFile   string
		jsdefinesMode bool
		definesMode   bool
		prettyprint   bool
	)
	flag.StringVar(&outputFile, "o", "", "Output file")
	flag.StringVar(&outputFile, "output", "", "Output file")
	flag.StringVar(&definesFile, "d", "", "defines.ace template")
	flag.StringVar(&definesFile, "defines", "", "defines.ace template")
	flag.BoolVar(&jsdefinesMode, "j", false, "Javascript defines mode")
	flag.BoolVar(&jsdefinesMode, "javascript", false, "Javascript defines mode")
	flag.BoolVar(&definesMode, "s", false, "Save the defines")
	flag.BoolVar(&definesMode, "savedefines", false, "Save the defines")
	flag.BoolVar(&prettyprint, "pp", true, "Pretty-print the output")
	flag.BoolVar(&prettyprint, "prettyprint", true, "Pretty-print the output")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage:\n  %s [options] [filename.ace]\n\nOptions:\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(64)
	}
	flag.Parse()

	inputFile := flag.Arg(0)
	if !definesMode && inputFile == "" {
		fmt.Fprintf(os.Stderr, "No template specified.\n")
		flag.Usage() // exits 64
		return
	}
	if definesMode && inputFile != "" {
		fmt.Fprintf(os.Stderr, "Extra template specified.\n")
		flag.Usage() // exits 64
		return
	}

	check := func(err error) {
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	chain := func(args ...interface{}) []interface{} { return args }
	dtmpl, err := templatetext.New(filepath.Base(definesFile)).Delims("[[", "]]").
		Funcs(templatetext.FuncMap{"Chain": chain}).ParseFiles(definesFile)
	check(err)
	dbuf := new(bytes.Buffer)
	check(dtmpl.Execute(dbuf, nil))

	aceopts := ace.InitializeOptions(&ace.Options{FuncMap: templatefunc.Funcs})
	defaultNoclose := aceopts.NoCloseTagNames
	definesAce, err := NewAceFile(definesFile, dbuf.Bytes(), aceopts)
	check(err)

	if !jsdefinesMode {
		indexAce, err := NewAceFile(inputFile, nil, aceopts)
		check(err)
		index, err := LoadAce(indexAce, &definesAce, Base(definesFile, aceopts), aceopts)
		check(err)
		text := Format(prettyprint, index.Tree.Root.String(), defaultNoclose)
		text += FormatSubtrees(prettyprint, index, defaultNoclose)
		check(WriteFile(outputFile, text))
		return
	}

	aceopts.NoCloseTagNames = []string{} // deviate from defaultNoclose
	aceopts.AttributeNameClass = "className"
	aceopts.FuncMap = templatefunc.JSXFuncs{}.MakeMap()

	defines, err := LoadAce(definesAce, nil, "", aceopts)
	check(err)

	if definesMode {
		check(WriteFile(outputFile, FormatSubtrees(prettyprint, defines, defaultNoclose)))
		return
	}

	jstemplate, err := templatetext.ParseFiles(inputFile)
	check(err)

	definesonly := templatetext.New("jsdefines").
		Funcs(templatetext.FuncMap(aceopts.FuncMap))

	definesTemplates := SortableTemplates(defines.Templates())
	sort.Stable(definesTemplates)

	for _, t := range definesTemplates {
		name := t.Name()
		if !strings.HasPrefix(name, definesAce.Basename+"::") {
			name = definesAce.Basename + name
		}
		tree := t.Tree
		if prettyprint {
			text := Format(true, tree.Root.String(), defaultNoclose)
			pretty, err := templatetext.New(name).
				Funcs(templatetext.FuncMap(aceopts.FuncMap)).Parse(text)
			check(err)
			tree = pretty.Tree
		}
		_, err := definesonly.AddParseTree(name, tree)
		check(err)
	}

	jdata := struct{ Defines []Define }{}
	for _, t := range definesTemplates {
		names := strings.Split(t.Name(), "::")
		tname := names[len(names)-1]
		if strings.HasPrefix(tname, "define_") {
			define, err := MakeDefine(definesonly, tname, definesAce.Basename+"::"+tname)
			check(err)
			jdata.Defines = append(jdata.Defines, define)
		}
	}
	buf := new(bytes.Buffer)
	check(jstemplate.Execute(buf, jdata))
	check(WriteFile(outputFile, buf.String()))
}

type Define struct {
	ShortName  string
	Iterable   string
	NeedList   bool
	UsesParams bool
	JSX        string
}

// SortableTemplates is for just sorting.
type SortableTemplates []*templatehtml.Template

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

	data := templatepipe.Data(Curly, t)
	if nota, ok := data.(templatepipe.Nota); ok {
		for k, v := range nota["Data"].(templatepipe.Nota) {
			if k == "Params" {
				define.UsesParams = true
			} else if k != "." {
				if define.Iterable != "" {
					return define, fmt.Errorf("Key %q is second: iterable already by %q", k, define.Iterable)
				}
				define.Iterable = k
				if n, ok := v.(templatepipe.Nota); ok {
					if _, ok := n["List"]; ok {
						define.NeedList = true
					}
				}
			}
		}
	}
	html := new(bytes.Buffer)
	if err := t.Execute(html, data); err != nil {
		return define, err
	}
	define.JSX = html.String()
	return define, nil
}

var vtype = reflect.TypeOf(templatepipe.Nota(nil))

func Curly(parent, key, full string) interface{} {
	if _, ok := vtype.MethodByName(key); ok {
		return nil
	}
	return templatepipe.CurlyX(parent, key, full)
}

// LoadAce is ace.Load without dealing with includes and setting Base'd names for the templates.
func LoadAce(base AceFile, inner *AceFile, innername string, opts *ace.Options) (*templatehtml.Template, error) {
	/*
		func LoadAce(basename, innername string, ..) {
		base, err := NewAceFile(basename, nil, opts)
		if err != nil {
			return nil, err
		}
		inner, err := NewAceFile(innername, nil, opts)
		if err != nil {
			return nil, err
		}
	*/
	if inner == nil {
		in, err := NewAceFile(innername, nil, opts)
		if err != nil {
			return nil, err
		}
		inner = &in
	}
	src := ace.NewSource(base.File, inner.File, nil)
	res, err := ace.ParseSource(src, opts)
	if err != nil {
		return nil, err
	}
	return ace.CompileResult(base.Basename, res, opts)
}

// Base returns filepath.Base'd filename sans extension if it matches opts.Extension.
func Base(filename string, opts *ace.Options) string {
	if filename == "" {
		return ""
	}
	n := filepath.Base(filename)
	if ext := filepath.Ext(n); ext == "."+opts.Extension {
		n = n[:len(n)-len(ext)]
	}
	return n
}

// NewAce returns *ace.File: may read filename if content is nil.
// Ace will know the file by base name.
func NewAceFile(filename string, content []byte, opts *ace.Options) (AceFile, error) {
	if content == nil && filename != "" {
		var err error
		content, err = ioutil.ReadFile(filename)
		if err != nil {
			return AceFile{}, err
		}
	}
	base := Base(filename, opts)
	return AceFile{
		File:     ace.NewFile(base, content),
		Basename: base,
	}, nil
}

type AceFile struct {
	*ace.File
	Basename string
}

// FormatSubtrees returns subtemplates trees forrmatted.
func FormatSubtrees(prettyprint bool, tpl *templatehtml.Template, noclose []string) (output string) {
	var names []string
	trees := map[string]*parse.Tree{}
	for _, x := range tpl.Templates() {
		name := x.Name()
		if name == tpl.Name() {
			continue // skip the root template
		}
		names = append(names, name)
		trees[name] = x.Tree
	}
	sort.Strings(names)
	for _, name := range names {
		output += fmt.Sprintf("{{/*\n*/}}{{define \"%s\"}}%s{{end}}", name, Format(prettyprint, trees[name].Root.String(), noclose))
	}
	return output
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

// Format pretty-formats html unless prettyprint if false.
func Format(prettyprint bool, input string, noclose []string) (output string) {
	if !prettyprint {
		// MAYBE html.Render() here
		return input
	}
	isnoclose := func(s string) bool {
		for _, x := range noclose {
			if x == s {
				return true
			}
		}
		return false
	}

	z := html.NewTokenizer(strings.NewReader(input))
	for level := 0; ; {
		tok := z.Next()
		tag, _ := z.TagName()
		raw := string(z.Raw())

		if tok == html.StartTagToken && !isnoclose(string(tag)) {
			level++
		} else if tok == html.EndTagToken && level != 0 {
			level--
		}

		if tok == html.DoctypeToken || tok == html.StartTagToken || tok == html.EndTagToken {
			output += raw[:len(raw)-1] + "\n" + strings.Repeat("  ", level) + ">"
		} else {
			if tok == html.TextToken {
				raw = strings.Trim(raw, "\n")
			}
			output += raw
		}
		if tok == html.ErrorToken {
			break
		}
	}
	return output
}
