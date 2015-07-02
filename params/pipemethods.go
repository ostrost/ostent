package params

import (
	"fmt"
	"html/template"
	"net/url"
)

func SprintfAttr(format string, args ...interface{}) template.HTMLAttr {
	return template.HTMLAttr(fmt.Sprintf(format, args...))
}

// FormActionAttr is for template.
func (q Query) FormActionAttr() interface{} {
	return SprintfAttr(" action=\"/form/%s\"", url.QueryEscape(q.ValuesEncode(nil)))
}

func (bp BoolParam) ToggleHrefAttr() interface{} {
	return SprintfAttr(" href=\"%s\"", bp.EncodeToggle())
}

func (pp PeriodParam) PeriodNameAttr() interface{} {
	return SprintfAttr(" name=%q", pp.Pname)
}

func (pp PeriodParam) PeriodValueAttr() interface{} {
	if pp.Input == "" {
		return template.HTMLAttr("")
	}
	return SprintfAttr(" value=\"%s\"", pp.Input)
}

func (pp PeriodParam) RefreshClassAttr(classes string) interface{} {
	if pp.InputErrd {
		classes += " has-warning"
	}
	return SprintfAttr(" class=%q", classes)
}

func (lp LimitParam) LessHrefAttr() interface{} {
	return SprintfAttr(" href=%q", lp.EncodeLess())
}

func (lp LimitParam) MoreHrefAttr() interface{} {
	return SprintfAttr(" href=%q", lp.EncodeMore())
}
