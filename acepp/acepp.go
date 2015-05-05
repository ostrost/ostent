package main

import (
	"flag"
	"fmt"
	templatehtml "html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	templatetext "text/template"

	"github.com/ostrost/ostent/acepp/templatep"
	"github.com/yosssi/ace"
)

func main() {
	var (
		outputFile  string
		definesFile string
		jscriptMode bool
		definesMode bool
	)
	// TODO MAYBE -prettyprint with "github.com/yosssi/gohtml" formatting
	flag.StringVar(&outputFile, "o", "", "Output file")
	flag.StringVar(&outputFile, "output", "", "Output file")
	flag.StringVar(&definesFile, "d", "", "defines.ace template")
	flag.StringVar(&definesFile, "defines", "", "defines.ace template")
	flag.BoolVar(&jscriptMode, "j", false, "Javascript mode")
	flag.BoolVar(&jscriptMode, "javascript", false, "Javascript mode")
	flag.BoolVar(&definesMode, "s", false, "Save the defines")
	flag.BoolVar(&definesMode, "savedefines", false, "Save the defines")
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
	aceopts := &ace.Options{
		DelimLeft:          "{{",    // default
		DelimRight:         "}}",    // default
		Extension:          "ace",   // default
		AttributeNameClass: "class", // default
		FuncMap:            templatep.AceFuncs,
	}

	if !jscriptMode {
		_, index, err := LoadAce(inputFile, definesFile, aceopts)
		check(err)
		trees := templatep.Subtrees(index)
		text := templatep.SprintfTrees(trees)
		text += index.Tree.Root.String()
		check(templatep.WriteFile(outputFile, text))
		return
	}

	aceopts.AttributeNameClass = "className"
	definesbase, defines, err := LoadAce(definesFile, "", aceopts)
	check(err)

	if definesMode {
		check(templatep.WriteTrees(outputFile, templatep.Subtrees(defines)))
		return
	}

	jscript, err := templatetext.
		New(Base(inputFile, aceopts)).
		Funcs(templatetext.FuncMap(templatep.AceFuncs)).
		ParseFiles(inputFile)
	check(err)

	for _, t := range defines.Templates() {
		_, err := jscript.AddParseTree(definesbase+t.Name(), t.Tree)
		check(err)
	}

	m := templatep.Data(&templatep.TextTemplate{Template: jscript}, jscriptMode)
	// jscriptMode is always true at this point

	s, err := templatep.StringExecute(jscript, m)
	check(err)

	// s = strings.Replace(s, "class=", "className=", -1)
	check(templatep.WriteFile(outputFile, s))
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
