package templatepipe

import (
	templatehtml "html/template"
	"strings"
	templatetext "text/template"
	"text/template/parse"
)

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

func Data(root Templater) interface{} {
	tree := Tree(root)
	if tree == nil {
		return "{}"
	}
	data := Dotted{}
	vars := map[string][]string{}
	for _, node := range tree.Root.Nodes {
		DataNode(root, node, &data, vars, nil)
	}
	return Mkmap(data, 0)
}

func DataNode(root Templater, node parse.Node, data *Dotted, vars map[string][]string, prefixwords []string) {
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
					data.Append(words, prefixwords)
				}
			}
			if withNode.List != nil {
				for _, n := range withNode.List.Nodes {
					DataNode(root, n, data, vars, prefixwords)
				}
			}
			if withNode.ElseList != nil {
				for _, n := range withNode.ElseList.Nodes {
					DataNode(root, n, data, vars, prefixwords)
				}
			}
			if withv != "" {
				delete(vars, withv)
			}
		case parse.NodeTemplate:
			tnode := node.(*parse.TemplateNode)
			var tawords []string
			for _, arg := range tnode.Pipe.Cmds[0].Args {
				s := arg.String()
				if len(s) > 1 && s[0] == '.' {
					tawords = strings.Split(s[1:], ".")
					data.Append(tawords, prefixwords)
					break // just one argument (pipeline) to "{{template}}" allowed anyway
				}
			}
			if lo := root.LookupT(tnode.Name); lo != nil {
				tr := Tree(lo)
				if tr != nil && tr.Root != nil {
					for _, n := range tr.Root.Nodes {
						pw := prefixwords
						if tawords != nil {
							pw = append(prefixwords, tawords...)
						}
						DataNode(root, n, data, vars, pw)
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
					case parse.NodeChain:
						chain := arg.(*parse.ChainNode)
						if chain.Node.Type() == parse.NodePipe {
							pipe := chain.Node.(*parse.PipeNode)
							for _, arg := range pipe.Cmds[0].Args {
								if arg.Type() == parse.NodeField {
									w := arg.String()
									if len(w) > 0 && w[0] == '.' {
										data.Append(strings.Split(w[1:], "."), prefixwords)
									}
								}
							}
						}
					case parse.NodeField:
						ident = arg.(*parse.FieldNode).Ident

						if len(decl) > 0 && len(decl[0].Ident) > 0 {
							vars[decl[0].Ident[0]] = ident
						}
						data.Append(ident, prefixwords)

					case parse.NodeVariable:
						ident = arg.(*parse.VariableNode).Ident

						if words, ok := vars[ident[0]]; ok {
							words := append(words, ident[1:]...)
							data.Append(words, prefixwords)
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
				case parse.NodeTemplate:
					// DataNode(root, ifnode, data, vars, prefixwords)
					arg0 := ifnode.(*parse.TemplateNode).Pipe.Cmds[0].Args[0]
					if arg0.Type() == parse.NodePipe {
						cmd0 := arg0.(*parse.PipeNode).Cmds[0]
						if cmd0.Type() == parse.NodeCommand {
							for _, a := range cmd0.Args {
								if s, prefix := a.String(), decl+"."; strings.HasPrefix(s, prefix) {
									keys = append(keys, strings.TrimPrefix(s, prefix))
								}
							}
						}
					}
				}
			}

			// fml
			arg0 := rangeNode.Pipe.Cmds[0].Args[0].String()
			words, ok := vars[arg0]
			if !ok && len(arg0) > 0 && arg0[0] == '.' {
				words = strings.Split(arg0[1:], ".")
				data.Append(words, prefixwords)
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
func (d *Dotted) Append(words []string, prefix []string) {
	if prefix != nil {
		words = append(prefix, words...)
	}
	if len(words) == 0 {
		return
	}
	if l := d.Leave(words[0]); l != nil {
		l.Append(words[1:], nil) // recursion
		return
	}
	n := &Dotted{Parent: d, Name: words[0]}
	d.Leaves = append(d.Leaves, n)
	n.Append(words[1:], nil) // recursion
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

func Curly(s string) string {
	if strings.HasSuffix(s, "HTML") {
		return /* "{" + */ "<span dangerouslySetInnerHTML={{__html: " + s + "}} />" // + ".props.children}"
	}
	return "{" + s + "}"
}

func Mkmap(top Dotted, level int) interface{} {
	if len(top.Leaves) == 0 {
		return Curly(top.Notation())
	}
	h := make(Hash)
	for _, l := range top.Leaves {
		if l.Ranged {
			if len(l.Keys) != 0 {
				kv := make(map[string]string)
				for _, k := range l.Keys {
					kv[k] = Curly(l.Decl + "." + k)
				}
				h[l.Name] = []map[string]string{kv}
			} else {
				h[l.Name] = []string{}
			}
		} else {
			h[l.Name] = Mkmap(*l, level+1)
		}
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
