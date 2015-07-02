package templatepipe

import (
	"strings"
	"text/template"
	"text/template/parse"
)

type CurlyFunc func(string, string, string) interface{}

func Data(cf CurlyFunc, root *template.Template) interface{} {
	if root == nil || root.Tree == nil || root.Tree.Root == nil {
		return "{}"
	}
	data := Dotted{}
	vars := map[string][]string{}
	for _, node := range root.Tree.Root.Nodes {
		DataNode(root, node, &data, vars, nil)
	}
	if cf == nil {
		cf = CurlyX
	}
	return Encurl(cf, data, 0)
}

func DataNode(root *template.Template, node parse.Node, data *Dotted, vars map[string][]string, prefixwords []string) {
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
			if sub := root.Lookup(tnode.Name); sub != nil && sub.Tree != nil && sub.Tree.Root != nil {
				for _, n := range sub.Tree.Root.Nodes {
					pw := prefixwords
					if tawords != nil {
						pw = append(prefixwords, tawords...)
					}
					DataNode(root, n, data, vars, pw)
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
									keys = append(keys, strings.Split(strings.TrimPrefix(s, prefix), ".")...)
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

// Nota is a type of encoded jsx notations for template transplitting.
type Nota map[string]interface{}

func (n Nota) String() string { return n.ToString() }

func (d *Dotted) Notation() (string, string, string) {
	if d == nil || d.Name == "" {
		return "", "", ""
	}
	if _, _, s := d.Parent.Notation(); s != "" {
		return s, d.Name, s + "." + d.Name
	}
	return "", d.Name, d.Name
}

func CurlyX(parent, key, full string) interface{} {
	return Curly(parent, key, full)
}

func Curly(parent, key, full string) string {
	return Curl(full)
}

func Curl(s string) string {
	if strings.HasSuffix(s, "HTML") {
		return /* "{" + */ "<span dangerouslySetInnerHTML={{__html: " + s + "}} />" // + ".props.children}"
	}
	return "{" + s + "}"
}

// Encurl returns constructed Nota.
func Encurl(cf CurlyFunc, parent Dotted, level int) interface{} {
	if len(parent.Leaves) == 0 {
		return cf(parent.Notation()) // may be nil
	}
	n := make(Nota)
	for _, l := range parent.Leaves {
		if l.Ranged {
			if len(l.Keys) != 0 {
				kv := make(map[string]Nota)
				for _, k := range l.Keys {
					if v := cf(l.Decl, k, l.Decl+"."+k); v != nil {
						kv[k] = make(Nota)
						kv[k]["."] = v
					}
				}
				if len(kv) != 0 {
					n[l.Name] = []map[string]Nota{kv}
				}
			} else {
				n[l.Name] = []string{}
			}
		} else if m := Encurl(cf, *l, level+1); m != nil {
			n[l.Name] = m
		}
	}
	_, _, n["."] = parent.Notation()
	return n
}
