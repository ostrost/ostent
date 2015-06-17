package templatep

import (
	"bytes"
	"encoding/json"
	"fmt"
	templatehtml "html/template"
	"net/url"
	"strings"
	templatetext "text/template"
	"text/template/parse"

	"github.com/ostrost/ostent/client"
)

// JS is whether we're doing it for jsx.
var JS bool

// classword returns either class or className depending on JS value.
func classword() string {
	return map[bool]string{
		false: "class",     // default
		true:  "className", // jsx case
	}[JS]
}

// forword returns either for or htmlFor depending on JS value.
func forword() string {
	return map[bool]string{
		false: "for",     // default
		true:  "htmlFor", // jsx case
	}[JS]
}

func CloseTagFunc(noclose []string) func(string) templatehtml.HTML {
	return func(tn string) templatehtml.HTML {
		for _, nc := range noclose {
			if tn == nc {
				return templatehtml.HTML("")
			}
		}
		return templatehtml.HTML("</" + tn + ">")
	}
}

// DotSplit splits s by last ".".
func DotSplit(s string) (string, string) {
	if s == "" {
		return "", ""
	}
	i := len(s) - 1
	for i > 0 && s[i] != '.' {
		i--
	}
	return s[:i], s[i+1:]
}

// DotSplitHash returns DotSplit of first (in no particular order) value from h.
func DotSplitHash(h Hash) (string, string) {
	var curled string
	for _, v := range h {
		curled = v.(string)
		break
		// First (no particular order) value is fine.
	}
	return DotSplit(uncurl(curled))
}

func toggleHrefAttr(value interface{}) interface{} {
	if JS {
		return fmt.Sprintf(" href={%s.Href} onClick={this.handleClick}", uncurl(value.(string)))
	}
	return templatehtml.HTMLAttr(fmt.Sprintf(" href=\"%s\"",
		value.(*client.BoolParam).EncodeToggle()))
}

func formActionAttr(query interface{}) interface{} {
	if JS {
		return fmt.Sprintf(" action={\"/form/\"+%s}", uncurl(query.(string)))
	}
	return templatehtml.HTMLAttr(fmt.Sprintf(" action=\"/form/%s\"",
		url.QueryEscape(query.(*client.Query).ValuesEncode(nil))))
}

func periodNameAttr(pparam interface{}) interface{} {
	if JS {
		prefix, _ := DotSplitHash(pparam.(Hash))
		_, pname := DotSplit(prefix)
		return fmt.Sprintf(" name=%q", pname)
	}
	period := pparam.(*client.PeriodParam)
	return templatehtml.HTMLAttr(fmt.Sprintf(" name=%q", period.Pname))
}

func periodValueAttr(pparam interface{}) interface{} {
	if JS {
		prefix, _ := DotSplitHash(pparam.(Hash))
		return fmt.Sprintf(" onChange={this.handleChange} value={%s.Input}", prefix)
	}
	if p := pparam.(*client.PeriodParam); p.Input != "" {
		return templatehtml.HTMLAttr(fmt.Sprintf(" value=\"%s\"", p.Input))
	}
	return templatehtml.HTMLAttr("")
}

/* TODO remove func refresh alltogether
func refresh(value interface{}) interface{} {
	if !JS {
		return value.(*client.Refresh)
	}
	prefix, _ := DotSplitHash(value.(Hash))
	// struct{struct{Duration, Above}}; Default} // mimic client.Refresh
	return struct {
		Period  string
		Default string
	}{
		Period:  fmt.Sprintf("{%s}", prefix),
		Default: fmt.Sprintf("{%s}", prefix),
		// TODO:
		// Period:    fmt.Sprintf("{%s.Period}", prefix),
		// Default:   fmt.Sprintf("{%s.Default}", prefix),
		// etc.
	}
}
// */

func ifDisabledAttr(value interface{}) templatehtml.HTMLAttr {
	if JS {
		return templatehtml.HTMLAttr(fmt.Sprintf("disabled={%s.Value ? \"disabled\" : \"\" }", uncurl(value.(string))))
	}
	if value.(*client.BoolParam).Value {
		return templatehtml.HTMLAttr("disabled=\"disabled\"")
	}
	return templatehtml.HTMLAttr("")
}

func ifClassAttr(value interface{}, classes ...string) (templatehtml.HTMLAttr, error) {
	s, err := ifClass(value, classes...)
	if err != nil {
		return templatehtml.HTMLAttr(""), err
	}
	if !JS {
		s = fmt.Sprintf("%q", s)
	}
	return templatehtml.HTMLAttr(fmt.Sprintf(" %s=%s", classword(), s)), nil
}

func ifClass(value interface{}, classes ...string) (string, error) {
	if len(classes) == 0 || len(classes) > 3 {
		return "", fmt.Errorf("number of args for ifClass*: either 2 or 3 or 4 got %d", 1+len(classes))
	}
	sndclass := ""
	if len(classes) > 1 {
		sndclass = classes[1]
	}
	fstclass := classes[0]
	if len(classes) > 2 {
		fstclass = classes[2] + " " + fstclass
		sndclass = classes[2] + " " + sndclass
	}
	if JS {
		return fmt.Sprintf("{%s.Value ? %q : %q }", uncurl(value.(string)), fstclass, sndclass), nil
	}
	if value.(*client.BoolParam).Value {
		return fstclass, nil
	}
	return sndclass, nil
}

func uncurl(s string) string {
	return strings.TrimSuffix(strings.TrimPrefix(s, "{"), "}")
}

func droplink(value interface{}, ss ...string) (interface{}, error) {
	if value == nil {
		return nil, fmt.Errorf("value supplied for droplink is nil")
	}
	var named string
	if len(ss) > 0 {
		named = ss[0]
	}
	AC := "text-right" // default
	if len(ss) > 1 {
		AC = ""
		if ss[1] != "" {
			AC = "text-" + ss[1]
		}
	}
	if JS {
		prefix, _ := DotSplitHash(value.(Hash))
		_, pname := DotSplit(prefix)
		enums := client.NewParamsENUM(nil)
		ed := enums[pname].EnumDecodec
		return client.DropLink{
			AlignClass: AC,
			Text:       ed.Text(named), // always static
			Href:       fmt.Sprintf("{%s.%s.%s}", prefix, named, "Href"),
			Class:      fmt.Sprintf("{%s.%s.%s}", prefix, named, "Class"),
			CaretClass: fmt.Sprintf("{%s.%s.%s}", prefix, named, "CaretClass"),
		}, nil
	}
	ep := value.(*client.EnumParam)
	pname, uptr := ep.EnumDecodec.Unew()
	if err := uptr.Unmarshal(named, new(bool)); err != nil {
		return nil, err
	}
	l := ep.EncodeUint(pname, uptr)
	l.AlignClass = AC
	return l, nil
}

func LabelClassColorPercent(p string) string {
	if len(p) > 2 { // 100% and more
		return "label label-danger"
	}
	if len(p) > 1 {
		if p[0] == '9' {
			return "label label-danger"
		}
		if p[0] == '8' {
			return "label label-warning"
		}
		if p[0] == '1' {
			return "label label-success"
		}
		return "label label-info"
	}
	return "label label-success"
}

func usepercent(val string) interface{} {
	var ca string
	if JS {
		ca = " className={LabelClassColorPercent(" + uncurl(val) + ")}"
	} else {
		ca = fmt.Sprintf(" class=%q", LabelClassColorPercent(val))
	}
	return struct {
		Value     string
		ClassAttr templatehtml.HTMLAttr
	}{
		Value:     val,
		ClassAttr: templatehtml.HTMLAttr(ca),
	}
}

func key(prefix, val string) templatehtml.HTMLAttr {
	if !JS {
		return templatehtml.HTMLAttr("")
	}
	return templatehtml.HTMLAttr(fmt.Sprintf(" key={%q+%s}", prefix+"-", uncurl(val)))
}

type Clipped struct {
	IDAttr      templatehtml.HTMLAttr
	ForAttr     templatehtml.HTMLAttr
	MWStyleAttr templatehtml.HTMLAttr
	Text        string
}

func clip(width int, prefix, val string, rest ...string) (*Clipped, error) {
	var key, mws string
	if JS {
		key = fmt.Sprintf("{%q+%s}", prefix+"-", uncurl(val))
		mws = fmt.Sprintf("{{maxWidth: '%dch'}}", width)
	} else { // quote everything
		key = fmt.Sprintf("%q", url.QueryEscape(prefix+"-"+val))
		mws = fmt.Sprintf("\"max-width: %dch \"", width)
	}
	if len(rest) == 1 {
		val = rest[0]
	} else if len(rest) > 0 {
		return nil, fmt.Errorf("clip expects either 5 or 6 arguments")
	}
	return &Clipped{
		IDAttr:      templatehtml.HTMLAttr("id=" + key),
		ForAttr:     templatehtml.HTMLAttr(forword() + "=" + key),
		MWStyleAttr: templatehtml.HTMLAttr("style=" + mws),
		Text:        val,
	}, nil
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

func curly(s string) string {
	if strings.HasSuffix(s, "HTML") {
		return /* "{" + */ "<span dangerouslySetInnerHTML={{__html: " + s + "}} />" // + ".props.children}"
	}
	return "{" + s + "}"
}

func mkmap(top Dotted, level int) interface{} {
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
			h[l.Name] = mkmap(*l, level+1)
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
	"dot":        dot,
	"key":        key,
	"clip":       clip,
	"droplink":   droplink,
	"usepercent": usepercent,

	"ifClass":         ifClass,
	"ifClassAttr":     ifClassAttr,
	"ifDisabledAttr":  ifDisabledAttr,
	"toggleHrefAttr":  toggleHrefAttr,
	"formActionAttr":  formActionAttr,
	"periodNameAttr":  periodNameAttr,
	"periodValueAttr": periodValueAttr,
	"closeTag":        CloseTagFunc(nil),
	"class":           classword,
	"for":             forword,

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
	return mkmap(data, 0)
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
