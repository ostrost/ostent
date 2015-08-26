// +build ignore

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
	dtmpl, err := templatetext.New(filepath.Base(definesFile)).Delims("[[", "]]").
		Funcs(templatetext.FuncMap{
		"Chain": func(args ...interface{}) []interface{} { return args },
	}).
		ParseFiles(definesFile)
	check(err)
	dbuf := new(bytes.Buffer)
	check(dtmpl.Execute(dbuf, nil))

	aceopts := ace.InitializeOptions(&ace.Options{FuncMap: templatefunc.Funcs})
	definesAce, err := NewAceFile(definesFile, dbuf.Bytes(), aceopts)
	check(err)

	if !jsdefinesMode {
		indexAce, err := NewAceFile(inputFile, nil, aceopts)
		check(err)
		index, err := LoadAce(indexAce, &definesAce, Base(definesFile, aceopts), aceopts)
		check(err)
		text := Format(prettyprint, index.Tree.Root.String(), aceopts.NoCloseTagNames)
		text += FormatSubtrees(prettyprint, index, aceopts.NoCloseTagNames)
		check(WriteFile(outputFile, text))
		return
	}

	noclose := aceopts.NoCloseTagNames
	aceopts.NoCloseTagNames = []string{}
	aceopts.AttributeNameClass = "className"
	aceopts.FuncMap = templatefunc.JSXFuncs{}.MakeMap()

	defines, err := LoadAce(definesAce, nil, "", aceopts)
	check(err)

	if definesMode {
		check(WriteFile(outputFile, FormatSubtrees(prettyprint, defines, noclose)))
		return
	}

	funcs := aceopts.FuncMap
	// override setrows & set the rowsset
	funcs["setrows"] = templatepipe.SetKFunc(".OverrideRows")
	funcs["rowsset"] = templatepipe.GetKFunc(".OverrideRows")

	jsdefines, err := templatetext.New(Base(inputFile, aceopts)).
		Funcs(templatetext.FuncMap(funcs)).ParseFiles(inputFile)
	check(err)

	for _, t := range defines.Templates() {
		name, tree := definesAce.Basename+t.Name(), t.Tree
		if prettyprint {
			text := Format(true, tree.Root.String(), noclose)
			y, err := templatetext.New(name).Funcs(templatetext.FuncMap(aceopts.FuncMap)).Parse(text)
			check(err)
			tree = y.Tree
		}
		_, err := jsdefines.AddParseTree(name, tree)
		check(err)
	}

	data := templatepipe.Data(Curly, jsdefines)
	buf := new(bytes.Buffer)
	check(jsdefines.Execute(buf, data))
	check(WriteFile(outputFile, buf.String()))
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
