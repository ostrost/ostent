package templatepipe

import (
	"fmt"
	"html/template"
	"strings"

	"github.com/ostrost/ostent/system/operating"
)

func Uncurl(s string) string {
	return strings.TrimSuffix(strings.TrimPrefix(s, "{"), "}")
}

func (n Nota) Uncurl() string {
	return Uncurl(n.ToString())
}

func (n Nota) FormActionAttr() interface{} {
	return fmt.Sprintf(" action={\"/form/\"+%s}", n.Uncurl())
}

func (n Nota) KeyAttr(prefix string) template.HTMLAttr {
	return template.HTMLAttr(fmt.Sprintf(" key={%q+%s}", prefix+"-", n.Uncurl()))
}

func (n Nota) Clip(width int, prefix string, id ...operating.ToStringer) (*operating.Clipped, error) {
	k, err := operating.ClipArgs(id, n.Uncurl())
	if err != nil {
		return nil, err
	}
	key := fmt.Sprintf("{%q+%s}", prefix+"-", Uncurl(k))
	return &operating.Clipped{
		IDAttr:      operating.SprintfAttr(" id=%s", key),
		ForAttr:     operating.SprintfAttr(" htmlFor=%s", key),
		MWStyleAttr: operating.SprintfAttr(" style={{maxWidth: '%dch'}}", width),
		Text:        n.ToString(),
	}, nil
}

func (n Nota) ToggleHrefAttr() interface{} {
	return fmt.Sprintf(" href={%s.Href} onClick={this.handleClick}", n.Uncurl())
}

func (n Nota) PeriodNameAttr() interface{} {
	_, pname := n.DotSplit()
	return fmt.Sprintf(" name=%q", pname)
}

func (n Nota) PeriodValueAttr() interface{} {
	return fmt.Sprintf(" value={%s.Input} onChange={this.handleChange}", n.ToString())
}

func (n Nota) RefreshClassAttr(classes string) interface{} {
	return fmt.Sprintf(" className={%q + (%s.InputErrd ? %q : \"\")}",
		classes, n.ToString(), " has-warning")
}

func (n Nota) LessHrefAttr() interface{} {
	return fmt.Sprintf(" href={%s.LessHref} onClick={this.handleClick}", n.Uncurl())
}

func (n Nota) MoreHrefAttr() interface{} {
	return fmt.Sprintf(" href={%s.MoreHref} onClick={this.handleClick}", n.Uncurl())
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

// Split calls DotSplit with n's string.
func (n Nota) DotSplit() (string, string) {
	return DotSplit(n.ToString())
}

func (n Nota) ToString() string {
	return n["."].(string)
}
