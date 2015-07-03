package operating

import (
	"fmt"
	"html/template"
	"net/url"
)

type Clipped struct {
	IDAttr      template.HTMLAttr
	ForAttr     template.HTMLAttr
	MWStyleAttr template.HTMLAttr
	Text        string
}

func Clip(width int, prefix string, id []fmt.Stringer, value string) (*Clipped, error) {
	key, err := ClipArgs(id, value)
	if err != nil {
		return nil, err
	}
	key = url.QueryEscape(prefix + "-" + key)
	return &Clipped{
		IDAttr:      SprintfAttr(" id=%q", key),
		ForAttr:     SprintfAttr(" for=%q", key),
		MWStyleAttr: SprintfAttr(" style=\"max-width: %dch \"", width),
		Text:        value,
	}, nil
}

func ClipArgs(id []fmt.Stringer, value string) (string, error) {
	if len(id) == 1 {
		return id[0].String(), nil
	} else if len(id) > 0 {
		return "", fmt.Errorf("Clip expects either 5 or 6 arguments")
	}
	return value, nil
}

func SprintfAttr(format string, args ...interface{}) template.HTMLAttr {
	return template.HTMLAttr(fmt.Sprintf(format, args...))
}
