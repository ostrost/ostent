package params

import (
	"fmt"
	"html/template"
	"net/url"
)

// FormActionAttr is for template.
func (q Query) FormActionAttr() (interface{}, error) {
	return template.HTMLAttr(fmt.Sprintf(" action=\"/form/%s\"",
		url.QueryEscape(q.ValuesEncode(nil)))), nil
}

func (bp BoolParam) ToggleHrefAttr() interface{} {
	return template.HTMLAttr(fmt.Sprintf(" href=\"%s\"", bp.EncodeToggle()))
}
