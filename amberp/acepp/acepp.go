package main

import (
	"flag"
	"fmt"
	templatehtml "html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	templatetext "text/template"

	"github.com/ostrost/ostent/amberp"
	"github.com/yosssi/ace"
)

func main() {
	var (
		outputFile  string
		definesFile string
		prettyPrint bool
		jscriptMode bool
		definesMode bool
	)
	flag.StringVar(&outputFile, "o", "", "Output file")
	flag.StringVar(&outputFile, "output", "", "Output file")
	flag.StringVar(&definesFile, "d", "", "Use defines file")
	flag.StringVar(&definesFile, "defines", "", "Use defines file")
	flag.BoolVar(&prettyPrint, "pp", false, "Pretty print output")
	flag.BoolVar(&prettyPrint, "prettyprint", false, "Pretty print output")
	flag.BoolVar(&jscriptMode, "j", false, "Javascript mode")
	flag.BoolVar(&jscriptMode, "javascript", false, "Javascript mode")
	flag.BoolVar(&definesMode, "s", false, "Save defines mode")
	flag.BoolVar(&definesMode, "savedefines", false, "Save defines mode")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage:\n  %s [options] [base.ace] [inner.ace]\n\nOptions:\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(64)
	}
	flag.Parse()

	inputFile := flag.Arg(0)
	if !definesMode && inputFile == "" {
		fmt.Fprintf(os.Stderr, "No input file specified.")
		flag.Usage()
		os.Exit(2) // TODO 64
	}

	check := func(err error) {
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	aceopt := &ace.Options{
		DelimLeft:  "{{",  // default
		DelimRight: "}}",  // default
		Extension:  "ace", // default
		FuncMap:    amberp.AceFuncs,
	}

	if !jscriptMode {
		_, index, err := LoadAce(inputFile, definesFile, aceopt)
		check(err)
		trees := amberp.Subtrees(index)
		text := amberp.SprintfTrees(trees)
		text += index.Tree.Root.String()
		check(amberp.WriteFile(outputFile, text))
		return
	}

	definesbase, defines, err := LoadAce(definesFile, "", aceopt)
	check(err)

	if definesMode {
		check(amberp.WriteTrees(outputFile, amberp.Subtrees(defines)))
		return
	}

	jscript, err := templatetext.New(filepath.Base(inputFile)).Funcs(templatetext.FuncMap(amberp.AceFuncs)).ParseFiles(inputFile)
	check(err)

	for _, t := range defines.Templates() {
		_, err := jscript.AddParseTree(definesbase+t.Name(), t.Tree)
		check(err)
	}

	m := amberp.Data(&amberp.TextTemplate{Template: jscript}, jscriptMode)

	s, err := amberp.StringExecute(jscript, m)
	check(err)

	s = strings.Replace(s, "class=", "className=", -1)
	check(amberp.WriteFile(outputFile, s))
}

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
