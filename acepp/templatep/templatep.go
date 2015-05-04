package templatep

import (
	"bytes"
	"encoding/json"
	"fmt"
	templatehtml "html/template"
	"io/ioutil"
	"os"
	"sort"
	"strings"
	templatetext "text/template"
	"text/template/parse"
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

type Hash map[string]interface{}

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
	h := make(Hash)
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
	return hindent(h.(Hash), 0)
}

func hindent(h Hash, level int) string {
	s := ""
	for k, v := range h {
		s += strings.Repeat(" ", level) + "(" + k + ")\n"
		vv, ok := v.(Hash)
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
	hashp *Hash
}

// func (dv dotValue) GoString() string { return dv.GoString() } // WTF?

func (dv dotValue) String() string {
	v := dv.s
	delete(*dv.hashp, "OVERRIDE")
	return v
}

func dot(v interface{}, key string) Hash {
	h := v.(Hash)
	h["OVERRIDE"] = dotValue{s: curly(key), hashp: &h}
	return h
}

// DotFuncs features "dot" function for templates. In use in acepp.
var DotFuncs = templatetext.FuncMap{"dot": dot}

// AceFuncs features functions for templates. In use in acepp and templates.
var AceFuncs = templatehtml.FuncMap{
	"dot": dot,
	"json": func(v interface{}) (string, error) {
		j, err := json.Marshal(v)
		return string(j), err
	},
}

type HTMLTemplate struct{ *templatehtml.Template }

func (ht HTMLTemplate) GetTree() *parse.Tree {
	if ht.Template == nil {
		return nil
	}
	return ht.Template.Tree
}
func (ht *HTMLTemplate) LookupT(n string) Templater { return &HTMLTemplate{ht.Lookup(n)} }

type TextTemplate struct{ *templatetext.Template }

func (tt TextTemplate) GetTree() *parse.Tree {
	if tt.Template == nil {
		return nil
	}
	return tt.Template.Tree
}
func (tt *TextTemplate) LookupT(n string) Templater { return &TextTemplate{tt.Lookup(n)} }

type Templater interface {
	GetTree() *parse.Tree
	LookupT(string) Templater
}

func Tree(root Templater) *parse.Tree {
	if root == nil {
		return nil
	}
	return root.GetTree()
}

// StringExecuteHTML does t.Execute into string returned. Does not clone.
func StringExecuteHTML(t *templatehtml.Template, data interface{}) (string, error) {
	buf := new(bytes.Buffer)
	if err := t.Execute(buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// StringExecute does t.Execute into string returned. Does not clone.
func StringExecute(t *templatetext.Template, data interface{}) (string, error) {
	buf := new(bytes.Buffer)
	if err := t.Execute(buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func Data(root Templater, jscriptMode bool) interface{} {
	tree := Tree(root)
	if tree == nil {
		return "{}"
	}
	data := Dotted{}
	vars := map[string][]string{}
	for _, node := range tree.Root.Nodes {
		DataNode(root, node, jscriptMode, &data, vars)
	}
	return mkmap(data, jscriptMode, 0)
}

func DataNode(root Templater, node parse.Node, jscriptMode bool, data *Dotted, vars map[string][]string) {
	if true {
		switch node.Type() {
		case parse.NodeWith:
			withNode := node.(*parse.WithNode)
			arg0 := withNode.Pipe.Cmds[0].Args[0].String()
			var withv string
			if len(arg0) > 0 && arg0[0] == '.' {
				if decl := withNode.Pipe.Decl; len(decl) > 0 {
					// just {{with $ := ...}} cases

					withv = decl[0].Ident[0]
					words := strings.Split(arg0[1:], ".")
					vars[withv] = words
					data.Append(words)
				}
			}
			if withNode.List != nil {
				for _, n := range withNode.List.Nodes {
					DataNode(root, n, jscriptMode, data, vars)
				}
			}
			if withNode.ElseList != nil {
				for _, n := range withNode.ElseList.Nodes {
					DataNode(root, n, jscriptMode, data, vars)
				}
			}
			if withv != "" {
				delete(vars, withv)
			}
		case parse.NodeTemplate:
			tnode := node.(*parse.TemplateNode)
			if lo := root.LookupT(tnode.Name); lo != nil {
				tr := Tree(lo)
				if tr != nil && tr.Root != nil {
					for _, n := range tr.Root.Nodes {
						DataNode(root, n, jscriptMode, data, vars)
					}
				}
			}
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
			words, ok := vars[arg0]
			if !ok && len(arg0) > 0 && arg0[0] == '.' {
				words = strings.Split(arg0[1:], ".")
				data.Append(words)
				ok = true
			}
			if ok {
				if leaf := data.Find(words); leaf != nil {
					leaf.Ranged = true
					leaf.Keys = append(leaf.Keys, keys...)
					leaf.Decl = decl // redefined $
				}
			}
		}
	}
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

func SaveDefines(outputfile, input string) error {
	T := struct {
		Name       string
		LeftDelim  string
		RightDelim string
	}{
		Name:       "zero",
		LeftDelim:  "[[",
		RightDelim: "]]",
	}
	// _ = template.New(T.Name).Funcs(DotFuncs).Delims(T.LeftDelim, T.RightDelim)
	trees, err := parse.Parse(T.Name, input, T.LeftDelim, T.RightDelim,
		DotFuncs, // .parseFuncs // template.FuncMap
		DotFuncs, // builtins // template.FuncMap
	)
	if err != nil {
		return err
	}
	delete(trees, T.Name)
	return WriteTrees(outputfile, trees)
}

func SprintfTrees(trees map[string]*parse.Tree) (text string) {
	for _, name := range KeysSorted(trees) {
		text += fmt.Sprintf("{{define \"%s\"}}%s{{end}}\n", name, trees[name].Root)
	}
	return text
}

func WriteTrees(outputfile string, trees map[string]*parse.Tree) error {
	return WriteFile(outputfile, SprintfTrees(trees))
}

// Subtrees returns a *parse.Tree map without tpl in it.
func Subtrees(tpl *templatehtml.Template) map[string]*parse.Tree {
	trees := map[string]*parse.Tree{}
	for _, x := range tpl.Templates() {
		if name := x.Name(); name != tpl.Name() {
			trees[name] = x.Tree
		}
	}
	return trees
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
