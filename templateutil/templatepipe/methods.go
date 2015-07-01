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

func (value Value) Uncurl() string {
	return Uncurl(value.ToString())
}

func (value Value) FormActionAttr() (interface{}, error) {
	return fmt.Sprintf(" action={\"/form/\"+%s}", value.Uncurl()), nil
}

func (value Value) KeyAttr(prefix string) template.HTMLAttr {
	return template.HTMLAttr(fmt.Sprintf(" key={%q+%s}", prefix+"-", value.Uncurl()))
}

func (value Value) Clip(width int, prefix string, id ...operating.ToStringer) (*operating.Clipped, error) {
	k, err := operating.ClipArgs(id, value.Uncurl())
	if err != nil {
		return nil, err
	}
	key := fmt.Sprintf("{%q+%s}", prefix+"-", Uncurl(k))
	return &operating.Clipped{
		IDAttr:      operating.SprintfAttr(" id=%s", key),
		ForAttr:     operating.SprintfAttr(" htmlFor=%s", key),
		MWStyleAttr: operating.SprintfAttr(" style={{maxWidth: '%dch'}}", width),
		Text:        value.ToString(),
	}, nil
}

func (value Value) ToggleHrefAttr() interface{} {
	return fmt.Sprintf(" href={%s.Href} onClick={this.handleClick}", value.Uncurl())
}

func (value Value) ToString() string { return string(value) }
func (value Value) forWord() string  { return "htmlFor" }
