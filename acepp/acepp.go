package main

import (
	"flag"
	"fmt"
	templatehtml "html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	templatetext "text/template"
	"text/template/parse"

	"code.google.com/p/go.net/html"
	"github.com/ostrost/ostent/acepp/templatep"
	"github.com/yosssi/ace"
)

func main() {
	var (
		outputFile  string
		definesFile string
		jscriptMode bool
		definesMode bool
		prettyprint bool
	)
	flag.StringVar(&outputFile, "o", "", "Output file")
	flag.StringVar(&outputFile, "output", "", "Output file")
	flag.StringVar(&definesFile, "d", "", "defines.ace template")
	flag.StringVar(&definesFile, "defines", "", "defines.ace template")
	flag.BoolVar(&jscriptMode, "j", false, "Javascript mode")
	flag.BoolVar(&jscriptMode, "javascript", false, "Javascript mode")
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
	aceopts := ace.InitializeOptions(&ace.Options{
		FuncMap: templatep.AceFuncs,
	})
	aceopts.FuncMap["closeTag"] = templatep.CloseTagFunc(aceopts.NoCloseTagNames)

	if !jscriptMode {
		_, index, err := LoadAce(inputFile, definesFile, aceopts)
		check(err)
		text := Format(prettyprint, index.Tree.Root.String(), aceopts.NoCloseTagNames)
		text += FormatSubtrees(prettyprint, index, aceopts)
		check(WriteFile(outputFile, text))
		return
	}

	aceopts.NoCloseTagNames = []string{}
	aceopts.AttributeNameClass = "className"
	aceopts.FuncMap["closeTag"] = templatep.CloseTagFunc(nil)
	templatep.JS = true

	definesbase, defines, err := LoadAce(definesFile, "", aceopts)
	check(err)

	if definesMode {
		check(WriteFile(outputFile, FormatSubtrees(prettyprint, defines, aceopts)))
		return
	}

	jscript, err := templatetext.
		New(Base(inputFile, aceopts)).
		Funcs(templatetext.FuncMap(templatep.AceFuncs)).
		ParseFiles(inputFile)
	check(err)

	for _, t := range defines.Templates() {
		name, tree := definesbase+t.Name(), t.Tree
		if prettyprint {
			text := Format(true, tree.Root.String(), aceopts.NoCloseTagNames)
			y, err := templatetext.New(name).Funcs(templatetext.FuncMap(aceopts.FuncMap)).Parse(text)
			check(err)
			tree = y.Tree
		}
		_, err := jscript.AddParseTree(name, tree)
		check(err)
	}

	m := templatep.Data(&templatep.TextTemplate{Template: jscript})

	s, err := templatep.StringExecute(jscript, m)
	check(err)

	check(WriteFile(outputFile, s))
}

// LoadAce is ace.Load without dealing with includes and setting Base'd names for the templates.
func LoadAce(basename, innername string, opts *ace.Options) (string, *templatehtml.Template, error) {
	base, err := ReadAce(basename, opts)
	if err != nil {
		return "", nil, err
	}
	inner, err := ReadAce(innername, opts)
	if err != nil {
		return "", nil, err
	}
	src := ace.NewSource(base, inner, nil)
	res, err := ace.ParseSource(src, opts)
	if err != nil {
		return "", nil, err
	}
	basebase := Base(basename, opts)
	template, err := ace.CompileResult(basebase, res, opts)
	return basebase, template, err
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

// ReadAce reads file and returns *ace.File Base-named.
func ReadAce(filename string, opts *ace.Options) (*ace.File, error) {
	var data []byte
	if filename != "" {
		var err error
		data, err = ioutil.ReadFile(filename)
		if err != nil {
			return nil, err
		}
	}
	return ace.NewFile(Base(filename, opts), data), nil
}

// FormatSubtrees returns subtemplates trees forrmatted.
func FormatSubtrees(prettyprint bool, tpl *templatehtml.Template, aceopts *ace.Options) (output string) {
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
		output += fmt.Sprintf("{{/*\n*/}}{{define \"%s\"}}%s{{end}}", name, Format(prettyprint, trees[name].Root.String(), aceopts.NoCloseTagNames))
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
