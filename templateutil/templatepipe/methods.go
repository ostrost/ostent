package templatepipe

import (
	"fmt"
	"strings"
)

func (value Value) Uncurl() string {
	return strings.TrimSuffix(strings.TrimPrefix(string(value), "{"), "}")
}

func (value Value) FormActionAttr() (interface{}, error) {
	return fmt.Sprintf(" action={\"/form/\"+%s}", value.Uncurl()), nil
}
