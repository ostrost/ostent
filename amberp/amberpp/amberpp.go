package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"text/template"
	"text/template/parse"

	"github.com/rzab/amber"
)

// Dotted is a tree.
type Dotted struct {
	Parent *Dotted
	Leaves []*Dotted // template/parse has nodes, we ought to have leaves
	Name   string

	Ranged bool
	Keys   []string
	Decl   string
}

// Append adds words into d.
func (d *Dotted) Append(words []string) {
	if len(words) == 0 {
		return
	}
	if l := d.Leave(words[0]); l != nil {
		l.Append(words[1:]) // recursion
		return
	}
	n := &Dotted{Parent: d, Name: words[0]}
	d.Leaves = append(d.Leaves, n)
	n.Append(words[1:]) // recursion
}

// Find traverses d to find by words.
func (d *Dotted) Find(words []string) *Dotted {
	if len(words) == 0 {
		return d
	}
	for _, l := range d.Leaves {
		if l.Name == words[0] {
			return l.Find(words[1:])
		}
	}
	return nil
}

func (d Dotted) Leave(name string) *Dotted {
	for _, l := range d.Leaves {
		if l.Name == name {
			return l
		}
	}
	return nil
}

func (d *Dotted) Notation() string {
	if d == nil || d.Name == "" {
		return ""
	}
	if s := d.Parent.Notation(); s != "" {
		return s + "." + d.Name
	}
	return d.Name
}

func (d Dotted) DebugString(level int) string {
	s := strings.Repeat(" ", level) + "[" + d.Name + "]\n"
	level += 2
	for _, l := range d.Leaves {
		s += l.DebugString(level)
	}
	level -= 2
	return s
}

func (d Dotted) GoString() string {
	return d.DebugString(0)
}

type hash map[string]interface{}

func curly(s string) string {
	if strings.HasSuffix(s, "HTML") {
		return /* "{" + */ "<span dangerouslySetInnerHTML={{__html: " + s + "}} />" // + ".props.children}"
	}
	return "{" + s + "}"
}

func mkmap(top Dotted, jscriptMode bool, level int) interface{} {
	if len(top.Leaves) == 0 {
		return curly(top.Notation())
	}
	h := make(hash)
	for _, l := range top.Leaves {
		if l.Ranged {
			if len(l.Keys) != 0 {
				kv := make(map[string]string)
				for _, k := range l.Keys {
					kv[k] = curly(l.Decl + "." + k)
				}
				h[l.Name] = []map[string]string{kv}
			} else {
				h[l.Name] = []string{}
			}
		} else {
			h[l.Name] = mkmap(*l, jscriptMode, level+1)
		}
	}
	if jscriptMode && level == 0 {
		h["CLASSNAME"] = "className"
	}
	return h
}

/* func string_hash(h interface{}) string {
	return hindent(h.(hash), 0)
}

func hindent(h hash, level int) string {
	s := ""
	for k, v := range h {
		s += strings.Repeat(" ", level) + "(" + k + ")\n"
		vv, ok := v.(hash)
		if ok && len(vv) > 0 {
			level += 2
			s += hindent(vv, level)
			level -= 2
		} else {
			s += strings.Repeat(" ", level + 2) + fmt.Sprint(v) + "\n"
		}
	}
	return s
} // */

type dotValue struct {
	s     string
	hashp *hash
}

func (dv dotValue) GoString() string {
	return dv.GoString()
}

func (dv dotValue) String() string {
	v := dv.s
	delete(*dv.hashp, "dot")
	return v
}

func dot(dot interface{}, key string) hash {
	h := dot.(hash)
	h["dot"] = dotValue{s: curly(key), hashp: &h}
	return h
}

var dotFuncs = map[string]interface{}{"dot": dot}

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

	fstplate, err := template.New("fst").Funcs(dotFuncs).Delims("[[", "]]").Parse(inputText)
	check(err)
	fst, err := StringExecute(fstplate, hash{})
	check(err)

	if !jscriptMode {
		check(writeFile(outputFile, fst))
		return
	}

	sndplate, err := template.New("snd").Funcs(template.FuncMap(amber.FuncMap)).Parse(fst)
	check(err)

	m := data(sndplate.Tree, jscriptMode)
	snd, err := StringExecute(sndplate, m)
	check(err)
	snd = regexp.MustCompile("</?script>").ReplaceAllLiteralString(snd, "")

	check(writeFile(outputFile, snd))
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
	// _ = template.New(T.Name).Funcs(dotFuncs).Delims(T.LeftDelim, T.RightDelim)
	trees, err := parse.Parse(T.Name, inputText, T.LeftDelim, T.RightDelim,
		dotFuncs, // .parseFuncs // template.FuncMap
		dotFuncs, // builtins // template.FuncMap
	)
	if err != nil {
		return err
	}
	var outputText string
	for name, t := range trees {
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

// StringExecute does t.Execute into string returned. Does not clone.
func StringExecute(t *template.Template, data interface{}) (string, error) {
	buf := new(bytes.Buffer)
	if err := t.Execute(buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func data(TREE *parse.Tree, jscriptMode bool) interface{} {
	if TREE == nil || TREE.Root == nil {
		return "{}" // mkmap(tree{})
	}

	data := Dotted{}
	vars := map[string][]string{}

	for _, node := range TREE.Root.Nodes { // here we go
		switch node.Type() {
		case parse.NodeAction:
			actionNode := node.(*parse.ActionNode)
			decl := actionNode.Pipe.Decl

			for _, cmd := range actionNode.Pipe.Cmds {
				if cmd.NodeType != parse.NodeCommand {
					continue
				}
				for _, arg := range cmd.Args {
					var ident []string
					switch arg.Type() {

					case parse.NodeField:
						ident = arg.(*parse.FieldNode).Ident

						if len(decl) > 0 && len(decl[0].Ident) > 0 {
							vars[decl[0].Ident[0]] = ident
						}
						data.Append(ident)

					case parse.NodeVariable:
						ident = arg.(*parse.VariableNode).Ident

						if words, ok := vars[ident[0]]; ok {
							words := append(words, ident[1:]...)
							data.Append(words)
							if len(decl) > 0 && len(decl[0].Ident) > 0 {
								vars[decl[0].Ident[0]] = words
							}
						}
					}
				}
			}
		case parse.NodeRange:
			rangeNode := node.(*parse.RangeNode)
			decl := rangeNode.Pipe.Decl[len(rangeNode.Pipe.Decl)-1].String()
			keys := []string{}

			for _, ifnode := range rangeNode.List.Nodes {
				switch ifnode.Type() {
				case parse.NodeAction:
					keys = append(keys, getKeys(decl, ifnode)...)
				case parse.NodeIf:
					for _, z := range ifnode.(*parse.IfNode).List.Nodes {
						if z.Type() == parse.NodeAction {
							keys = append(keys, getKeys(decl, z)...)
						}
					}
				}
			}

			// fml
			arg0 := rangeNode.Pipe.Cmds[0].Args[0].String()
			if words, ok := vars[arg0]; ok {
				if leaf := data.Find(words); leaf != nil {
					leaf.Ranged = true
					leaf.Keys = append(leaf.Keys, keys...)
					leaf.Decl = decl // redefined $
				}
			}
		}
	}
	return mkmap(data, jscriptMode, 0)
}

func getKeys(decl string, parseNode parse.Node) (keys []string) {
	for _, cmd := range parseNode.(*parse.ActionNode).Pipe.Cmds {
		if cmd.NodeType != parse.NodeCommand {
			continue
		}
		for _, arg := range cmd.Args {
			if arg.Type() != parse.NodeVariable {
				continue
			}
			ident := arg.(*parse.VariableNode).Ident
			if len(ident) < 2 || ident[0] != decl {
				continue
			}
			keys = append(keys, ident[1])
		}
	}
	return
}
