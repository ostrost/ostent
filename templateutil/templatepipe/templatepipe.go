package templatepipe

import (
	"fmt"
	"strings"
	"text/template"
	"text/template/parse"
)

/*
func (n Nota) uncurl() string   { return uncurl(n.String()) }
func (c curled) uncurl() string { return uncurl(string(c)) }
func uncurl(s string) string {
	return strings.TrimSuffix(strings.TrimPrefix(s, "{"), "}")
} // */

type curlyFunc func(string, string, string) interface{}

func templateData(cf curlyFunc, root *template.Template) interface{} {
	if root == nil || root.Tree == nil || root.Tree.Root == nil {
		return "{}"
	}
	c := context{vars: make(map[string][]string)}
	for _, n := range root.Tree.Root.Nodes {
		c.node(n, nodeArgs{template: root})
	}
	if cf == nil {
		cf = curly
	}
	return encurl(cf, c.dotted, 0)
}

type nodeArgs struct {
	template    *template.Template
	prefixWords []string
	decl        []*parse.VariableNode
}

type context struct {
	dotted dotted              // end result
	vars   map[string][]string // for vars kept between Node* methods calls
}

func (c *context) touch(na nodeArgs, words []string) {
	c.dotted.append(words, na.prefixWords)
	if len(na.decl) > 0 && len(na.decl[0].Ident) > 0 {
		c.vars[na.decl[0].Ident[0]] = words
	}
}

func (c *context) ranging(rangeNode *parse.RangeNode, na nodeArgs) {
	var (
		decl  = rangeNode.Pipe.Decl[len(rangeNode.Pipe.Decl)-1].String()
		field = rangeNode.Pipe.Cmds[0].Args[0].(*parse.FieldNode)
	)
	leaf := c.dotted.append(field.Ident, na.prefixWords)
	leaf.keys = append(leaf.keys, getKeys(decl, rangeNode.List)...)
	leaf.decl = decl // redefined $
	leaf.ranged = true
}

// nolint: gocyclo
func (c *context) node(node parse.Node, na nodeArgs) {
	switch node.Type() {
	// straightforward recursives:
	case parse.NodeCommand:
		for _, n := range node.(*parse.CommandNode).Args {
			c.node(n, na)
		}
	case parse.NodeIf:
		fi := node.(*parse.IfNode)
		if fi.List != nil {
			c.node(fi.List, na)
		}
		if fi.ElseList != nil {
			c.node(fi.ElseList, na)
		}
	case parse.NodeList:
		for _, n := range node.(*parse.ListNode).Nodes {
			c.node(n, na)
		}
	case parse.NodePipe:
		for _, n := range node.(*parse.PipeNode).Cmds {
			c.node(n, na)
		}
	case parse.NodeWith:
		with := node.(*parse.WithNode)
		c.node(with.Pipe, na)
		if with.List != nil {
			c.node(with.List, na)
		}
		if with.ElseList != nil {
			c.node(with.ElseList, na)
		}

	// other recursives:
	case parse.NodeAction:
		an := node.(*parse.ActionNode)
		na.decl = an.Pipe.Decl // !
		c.node(an.Pipe, na)
	case parse.NodeTemplate:
		t := node.(*parse.TemplateNode)
		c.node(t.Pipe, na)
		if s := na.template.Lookup(t.Name); s != nil && s.Tree != nil && s.Tree.Root != nil {
			c.node(s.Tree.Root, na)
		}

	// touchers:
	case parse.NodeRange:
		c.ranging(node.(*parse.RangeNode), na)
	case parse.NodeField:
		c.touch(na, node.(*parse.FieldNode).Ident)
	case parse.NodeVariable:
		v := node.(*parse.VariableNode)
		if words, ok := c.vars[v.Ident[0]]; ok {
			c.touch(na, append(words, v.Ident[1:]...))
		}
	}
}

// nolint: gocyclo
func getKeys(decl string, node parse.Node) (keys []string) {
	switch node.Type() {
	case parse.NodeAction:
		return getKeys(decl, node.(*parse.ActionNode).Pipe)
	case parse.NodeIf:
		return getKeys(decl, node.(*parse.IfNode).List)
	case parse.NodeTemplate:
		return getKeys(decl, node.(*parse.TemplateNode).Pipe)
	case parse.NodeCommand:
		for _, n := range node.(*parse.CommandNode).Args {
			keys = append(keys, getKeys(decl, n)...)
		}
	case parse.NodeList:
		for _, n := range node.(*parse.ListNode).Nodes {
			keys = append(keys, getKeys(decl, n)...)
		}
	case parse.NodePipe:
		for _, n := range node.(*parse.PipeNode).Cmds {
			keys = append(keys, getKeys(decl, n)...)
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

// dotted is a tree.
type dotted struct {
	parent *dotted
	leaves []*dotted // template/parse has nodes, we ought to have leaves
	name   string

	ranged bool
	keys   []string
	decl   string
}

// append adds words into d.
func (d *dotted) append(words []string, prefix []string) *dotted {
	if prefix != nil {
		words = append(prefix, words...)
	}
	if len(words) == 0 {
		return d
	}
	if l := d.leave(words[0]); l != nil {
		return l.append(words[1:], nil)
	}
	n := &dotted{parent: d, name: words[0]}
	d.leaves = append(d.leaves, n)
	return n.append(words[1:], nil)
}

// find traverses d to search by words. Used by tests.
func (d *dotted) find(words []string) *dotted {
	if len(words) == 0 {
		return d
	}
	for _, l := range d.leaves {
		if l.name == words[0] {
			return l.find(words[1:])
		}
	}
	return nil
}

func (d dotted) leave(name string) *dotted {
	for _, l := range d.leaves {
		if l.name == name {
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
	return string(v.(curled))
}

func (d *dotted) notation() (string, string, string) {
	if d == nil || d.name == "" {
		return "", "", ""
	}
	if _, _, s := d.parent.notation(); s != "" {
		return s, d.name, s + "." + d.name
	}
	return "", d.name, d.name
}

func curly(_, _, full string) interface{} { return curled(curl(full)) }

type curled string

func curl(s string) string {
	if strings.HasSuffix(s, "HTML") {
		return /* "{" + */ "<span dangerouslySetInnerHTML={{__html: " + s + "}} />" // + ".props.children}"
	}
	return "{" + s + "}"
}

// encurl returns constructed Nota.
func encurl(cf curlyFunc, parent dotted, level int) interface{} {
	if len(parent.leaves) == 0 {
		return cf(parent.notation()) // may be nil
	}
	n := make(Nota)
	for _, l := range parent.leaves {
		if l.ranged {
			if len(l.keys) != 0 {
				kv := make(map[string]Nota)
				for _, k := range l.keys {
					if v := cf(l.decl, k, l.decl+"."+k); v != nil {
						kv[k] = make(Nota)
						kv[k]["."] = v
					}
				}
				if len(kv) != 0 {
					n[l.name] = []map[string]Nota{kv}
				}
			} else {
				n[l.name] = []string{}
			}
		} else if m := encurl(cf, *l, level+1); m != nil {
			n[l.name] = m
		}
	}
	_, _, n["."] = parent.notation()
	return n
}
