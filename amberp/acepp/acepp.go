package main

import (
	"flag"
	"fmt"
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

	d := definesFile
	definesprefix := d[:len(d)-len(filepath.Ext(d))]

	if !jscriptMode {
		d := inputFile
		inputprefix := d[:len(d)-len(filepath.Ext(d))]
		index, err := ace.Load(inputprefix, definesprefix, &ace.Options{FuncMap: amberp.AceFuncs})
		check(err)
		trees := amberp.Subtrees(index)
		treeprefix := filepath.Dir(definesprefix) + string(os.PathSeparator)
		text := amberp.SprintfPrefixedTrees(treeprefix, trees)
		text += amberp.SprintfTrees(trees)
		text += index.Tree.Root.String()
		check(amberp.WriteFile(outputFile, text))
		return
	}

	defines, err := ace.Load(definesprefix, "", &ace.Options{FuncMap: amberp.AceFuncs})
	check(err)

	if definesMode {
		check(amberp.WriteTrees(outputFile, amberp.Subtrees(defines)))
		return
	}

	jscript, err := templatetext.New(filepath.Base(inputFile)).Funcs(templatetext.FuncMap(amberp.AceFuncs)).ParseFiles(inputFile)
	check(err)

	for _, t := range defines.Templates() {
		_, err := jscript.AddParseTree(filepath.Base(definesprefix)+t.Name(), t.Tree)
		check(err)
	}

	m := amberp.Data(&amberp.TextTemplate{Template: jscript}, jscriptMode)

	s, err := amberp.StringExecute(jscript, m)
	check(err)

	s = strings.Replace(s, "class=", "className=", -1)
	check(amberp.WriteFile(outputFile, s))
}
