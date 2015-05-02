package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"sort"
	"text/template"
	"text/template/parse"

	"github.com/ostrost/ostent/amberp"
	"github.com/rzab/amber"
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
	flag.Parse()

	inputFile := flag.Arg(0)
	if !definesMode && inputFile == "" {
		fmt.Fprintf(os.Stderr, "No input file specified.")
		flag.Usage()
		os.Exit(2)
	}

	check := func(err error) {
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	inputText := ""
	if definesFile != "" {
		b, err := ioutil.ReadFile(definesFile)
		check(err)
		newText, err := compile(b, prettyPrint, jscriptMode)
		check(err)
		inputText += newText
		if inputText[len(inputText)-1] == '\n' { // amber does add this '\n', which is fine for the end of a file, which inputText is not
			inputText = inputText[:len(inputText)-1]
		}
	}

	if definesMode {
		check(saveDefines(outputFile, inputText))
		return
	}

	b, err := ioutil.ReadFile(inputFile)
	check(err)
	newText, err := compile(b, prettyPrint, jscriptMode)
	check(err)
	inputText += newText

	fstplate, err := template.New("fst").Funcs(amberp.DotFuncs).Delims("[[", "]]").Parse(inputText)
	check(err)
	fst, err := amberp.StringExecute(fstplate, amberp.Hash{})
	check(err)

	if !jscriptMode {
		check(writeFile(outputFile, fst))
		return
	}

	sndplate, err := template.New("snd").Funcs(template.FuncMap(amber.FuncMap)).Parse(fst)
	check(err)

	m := amberp.Data(&amberp.TextTemplate{Template: sndplate}, jscriptMode)
	snd, err := amberp.StringExecute(sndplate, m)
	check(err)
	snd = regexp.MustCompile("</?script>").ReplaceAllLiteralString(snd, "")

	check(writeFile(outputFile, snd))
}

func KeysSorted(trees map[string]*parse.Tree) []string {
	keys := make([]string, len(trees))
	i := 0
	for k := range trees {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	return keys
}

func saveDefines(outputFile, inputText string) error {
	T := struct {
		Name       string
		LeftDelim  string
		RightDelim string
	}{
		Name:       "zero",
		LeftDelim:  "[[",
		RightDelim: "]]",
	}
	// _ = template.New(T.Name).Funcs(amberp.DotFuncs).Delims(T.LeftDelim, T.RightDelim)
	trees, err := parse.Parse(T.Name, inputText, T.LeftDelim, T.RightDelim,
		amberp.DotFuncs, // .parseFuncs // template.FuncMap
		amberp.DotFuncs, // builtins // template.FuncMap
	)
	if err != nil {
		return err
	}
	var outputText string
	for _, name := range KeysSorted(trees) {
		t := trees[name]
		if name == T.Name { // skip the toplevel
			continue
		}
		if t == nil || t.Root == nil {
			continue
		}
		outputText += fmt.Sprintf("{{define \"%s\"}}%s{{end}}\n", name, t.Root)
	}
	return writeFile(outputFile, outputText)
}

func writeFile(optFilename, s string) error {
	b := []byte(s)
	if optFilename != "" {
		return ioutil.WriteFile(optFilename, b, 0644)
	}
	_, err := os.Stdout.Write(b)
	return err
}

func compile(input []byte, prettyPrint, jscriptMode bool) (string, error) {
	compiler := amber.New()
	compiler.PrettyPrint = prettyPrint
	if jscriptMode {
		compiler.ClassName = "className"
	}
	if err := compiler.Parse(string(input)); err != nil {
		return "", err
	}
	return compiler.CompileString()
}
