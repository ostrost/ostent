package client

import (
	"fmt"
	"strings"
)

// Upointer dictated methods:

// Once there're no renamed parameters, these Unmarshal should do.
// func (r *UintDF) Unmarshal(s string) error { return UnmarshalStringFunc(r.UnmarshalJSON)(s) }
// func (r *UintPS) Unmarshal(s string) error { return UnmarshalStringFunc(r.UnmarshalJSON)(s) }

// Unmarshal for UintDF. Knows renamed parameter.
func (r *UintDF) Unmarshal(data string) error {
	issize := data == "size"
	if issize {
		data = "dfsize"
	}
	return UnmarshalMaybe(issize, r.UnmarshalJSON, data)
}

// Unmarshal for UintPS. Knows renamed parameter.
func (r *UintPS) Unmarshal(data string) error {
	issize := data == "size"
	if issize {
		data = "pssize"
	}
	isuser := data == "user"
	if isuser {
		data = "uid"
	}
	return UnmarshalMaybe(issize || isuser, r.UnmarshalJSON, data)
}

// Uinter dictated methods:

func (r UintDF) Touint() Uint             { return Uint(r) }
func (r UintDF) Marshal() (string, error) { return MarshalStringFunc(r.MarshalJSON)() }

func (r UintPS) Touint() Uint             { return Uint(r) }
func (r UintPS) Marshal() (string, error) { return MarshalStringFunc(r.MarshalJSON)() }

// Helpers:

func UnmarshalMaybe(rename bool, unmarshal BytesUnmarshal, data string) error {
	if err := UnmarshalStringFunc(unmarshal)(data); err != nil || !rename {
		return err
	}
	return RenamedConstError("")
}

func UnmarshalStringFunc(unmarshaler BytesUnmarshal) func(string) error {
	return func(data string) error {
		return unmarshaler([]byte(fmt.Sprintf("%q", strings.ToUpper(data))))
	}
}

func MarshalStringFunc(marshaler BytesEnmarshal) func() (string, error) {
	return func() (string, error) {
		b, err := marshaler()
		if err != nil {
			return "", err
		}
		if l := len(b); l > 2 && b[0] == '"' && b[l-1] == '"' {
			b = b[1 : l-1]
		}
		s := strings.ToLower(string(b))
		return s, nil
	}
}

type BytesEnmarshal func() ([]byte, error)
type BytesUnmarshal func([]byte) error
