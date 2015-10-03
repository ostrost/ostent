package templatepipe

import (
	"fmt"
	"strings"
	"text/template"
	"text/template/parse"
)

// Uncurl is required by templatefunc.Uncurler interface.
func (n Nota) Uncurl() string   { return Uncurl(n.String()) }
func (c Curled) Uncurl() string { return Uncurl(string(c)) }

func Uncurl(s string) string {
	return strings.TrimSuffix(strings.TrimPrefix(s, "{"), "}")
}

type CurlyFunc func(string, string, string) interface{}

func Data(cf CurlyFunc, root *template.Template) interface{} {
	if root == nil || root.Tree == nil || root.Tree.Root == nil {
		return "{}"
	}
	c := Context{Vars: make(map[string][]string)}
	for _, node := range root.Tree.Root.Nodes {
		c.Node(node, NodeArgs{Template: root})
	}
	if cf == nil {
		cf = CurlyX
	}
	return Encurl(cf, c.Dotted, 0)
}

type NodeArgs struct {
	Template    *template.Template
	PrefixWords []string
	Decl        []*parse.VariableNode
}

type Context struct {
	Dotted Dotted              // end result
	Vars   map[string][]string // for vars kept between Node* methods calls
}

func (c *Context) Touch(na NodeArgs, words []string) {
	c.Dotted.Append(words, na.PrefixWords)
	if len(na.Decl) > 0 && len(na.Decl[0].Ident) > 0 {
		c.Vars[na.Decl[0].Ident[0]] = words
	}
}

func (c *Context) Ranging(rangeNode *parse.RangeNode, na NodeArgs) {
	var (
		decl  = rangeNode.Pipe.Decl[len(rangeNode.Pipe.Decl)-1].String()
		field = rangeNode.Pipe.Cmds[0].Args[0].(*parse.FieldNode)
	)
	leaf := c.Dotted.Append(field.Ident, na.PrefixWords)
	leaf.Keys = append(leaf.Keys, getKeys(decl, rangeNode.List)...)
	leaf.Decl = decl // redefined $
	leaf.Ranged = true
}

func (c *Context) Node(node parse.Node, na NodeArgs) {
	switch node.Type() {
	// plain recursives:
	case parse.NodeCommand:
		for _, n := range node.(*parse.CommandNode).Args {
			c.Node(n, na)
		}
	case parse.NodeList:
		for _, n := range node.(*parse.ListNode).Nodes {
			c.Node(n, na)
		}
	case parse.NodePipe:
		for _, n := range node.(*parse.PipeNode).Cmds {
			c.Node(n, na)
		}

	// recursives:
	case parse.NodeAction:
		an := node.(*parse.ActionNode)
		na.Decl = an.Pipe.Decl // !
		c.Node(an.Pipe, na)
	case parse.NodeWith:
		with := node.(*parse.WithNode)
		c.Node(with.Pipe, na)
		c.Node(with.List, na)
		c.Node(with.ElseList, na)
	case parse.NodeTemplate:
		t := node.(*parse.TemplateNode)
		c.Node(t.Pipe, na)
		if s := na.Template.Lookup(t.Name); s != nil && s.Tree != nil && s.Tree.Root != nil {
			c.Node(s.Tree.Root, na)
		}

	// touchers:
	case parse.NodeRange:
		c.Ranging(node.(*parse.RangeNode), na)
	case parse.NodeField:
		c.Touch(na, node.(*parse.FieldNode).Ident)
	case parse.NodeVariable:
		v := node.(*parse.VariableNode)
		if words, ok := c.Vars[v.Ident[0]]; ok {
			c.Touch(na, append(words, v.Ident[1:]...))
		}
	}
}

func getKeys(decl string, node parse.Node) (keys []string) {
	switch node.Type() {
	case parse.NodeAction:
		return getKeys(decl, node.(*parse.ActionNode).Pipe)
	case parse.NodeIf:
		return getKeys(decl, node.(*parse.IfNode).List)
	case parse.NodeTemplate:
		return getKeys(decl, node.(*parse.TemplateNode).Pipe)
	case parse.NodeCommand:
		for _, arg := range node.(*parse.CommandNode).Args {
			keys = append(keys, getKeys(decl, arg)...)
		}
	case parse.NodeList:
		for _, n := range node.(*parse.ListNode).Nodes {
			keys = append(keys, getKeys(decl, n)...)
		}
	case parse.NodePipe:
		for _, cmd := range node.(*parse.PipeNode).Cmds {
			keys = append(keys, getKeys(decl, cmd)...)
		}
	case parse.NodeVariable:
		ident := node.(*parse.VariableNode).Ident
		if len(ident) < 2 || ident[0] != decl {
			panic(fmt.Errorf("Unexpected ident: %+v", ident))
		}
		return ident[1:]
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
func (d *Dotted) Append(words []string, prefix []string) *Dotted {
	if prefix != nil {
		words = append(prefix, words...)
	}
	if len(words) == 0 {
		return d
	}
	if l := d.Leave(words[0]); l != nil {
		return l.Append(words[1:], nil)
	}
	n := &Dotted{Parent: d, Name: words[0]}
	d.Leaves = append(d.Leaves, n)
	return n.Append(words[1:], nil)
}

// Find traverses d to search by words. Used by tests.
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

func (n Nota) String() string {
	v := n["."]
	if s, ok := v.(string); ok {
		return s
	}
	return string(v.(Curled))
}

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

type Curled string

func Curly(parent, key, full string) Curled {
	return Curled(Curl(full))
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
