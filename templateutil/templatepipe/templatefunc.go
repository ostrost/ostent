package templatepipe

import "strings"

// SetKFunc constructs a func which
// sets k key to Curly(string (n))
// in passed interface{} (v) being a Nota.
// SetKFunc is used by acepp only.
func SetKFunc(k string) func(interface{}, string) interface{} {
	return func(v interface{}, n string) interface{} {
		if args := strings.Split(n, " "); len(args) > 1 {
			var list []string
			for _, arg := range args {
				list = append(list, Curl(arg))
			}
			v.(Nota)[k] = list
			return v
		}
		v.(Nota)[k] = Curl(n)
		return v
	}
}

// GetKFunc constructs a func which
// gets, deletes and returns k key
// in passed interface{} (v) being a Nota.
// GetKFunc is used by acepp only.
func GetKFunc(k string) func(interface{}) interface{} {
	return func(v interface{}) interface{} {
		h, ok := v.(Nota)
		if !ok {
			return "" // empty pipeline, affects dispatch
		}
		n := h[k]
		if args, ok := n.([]string); ok {
			if len(args) > 1 {
				h[k] = args[1:]
			} else {
				delete(h, k)
			}
			return args[0]
		}
		delete(h, k)
		return n // may also be empty, affects dispatch
	}
}
