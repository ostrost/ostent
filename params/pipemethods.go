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

func (bp BoolParam) BoolParamClassAttr(classes ...string) (template.HTMLAttr, error) {
	fstclass, sndclass, err := ClassesChoices("BoolParamClassAttr", classes)
	if err != nil {
		return template.HTMLAttr(""), err
	}
	s := fstclass
	if !bp.Value {
		s = sndclass
	}
	return SprintfAttr(" class=%q", s), nil
}

// TODO dup from operating
func ClassesChoices(caller string, classes []string) (string, string, error) {
	if len(classes) == 0 || len(classes) > 3 {
		return "", "", fmt.Errorf("number of args for %s: either 2 or 3 or 4 got %d", caller, 1+len(classes))
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
	return fstclass, sndclass, nil
}

func (bp BoolParam) DisabledAttr() interface{} {
	if !bp.Value {
		return template.HTMLAttr("")
	}
	return SprintfAttr(" disabled=%q", "disabled")
}

func (bp BoolParam) ToggleHrefAttr() interface{} {
	return SprintfAttr(" href=\"%s\"", bp.EncodeToggle())
}

func (ep EnumParam) EnumClassAttr(named, classif string, optelse ...string) (template.HTMLAttr, error) {
	classelse, err := EnumClassAttrArgs(optelse)
	if err != nil {
		return template.HTMLAttr(""), err
	}
	_, uptr := ep.EnumDecodec.Unew()
	if err := uptr.Unmarshal(named, new(bool)); err != nil {
		return template.HTMLAttr(""), err
	}
	if ep.Number.Uint != uptr.Touint() {
		classif = classelse
	}
	return SprintfAttr(" class=%q", classif), nil
}

func EnumClassAttrArgs(opt []string) (string, error) {
	if len(opt) == 1 {
		return opt[0], nil
	} else if len(opt) > 1 {
		return "", fmt.Errorf("number of args for EnumClassAttr: either 2 or 3 got %d",
			2+len(opt))
	}
	return "", nil
}

func (ep EnumParam) EnumLink(args ...string) (EnumLink, error) {
	named, aclass := EnumLinkArgs(args)
	pname, uptr := ep.EnumDecodec.Unew()
	if err := uptr.Unmarshal(named, new(bool)); err != nil {
		return EnumLink{}, err
	}
	l := ep.EncodeUint(pname, uptr)
	l.AlignClass = aclass
	return l, nil
}

func EnumLinkArgs(args []string) (string, string) {
	var named string
	if len(args) > 0 {
		named = args[0]
	}
	aclass := "text-right" // default
	if len(args) > 1 {
		aclass = ""
		if args[1] != "" {
			aclass = "text-" + args[1]
		}
	}
	return named, aclass
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
