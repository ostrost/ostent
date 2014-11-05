// +build darwin

package types

import "fmt"

func RAMFields(ram RAM) []NameString {
	return []NameString{
		{"memory-free", fmt.Sprintf("%d", ram.Raw.Free)},
		{"memory-inactive", fmt.Sprintf("%d", ram.Raw.ActualFree-ram.Raw.Free)},
		{"memory-wired", fmt.Sprintf("%d", ram.Extra1)},
		{"memory-active", fmt.Sprintf("%d", ram.Extra2)},
	}
}
