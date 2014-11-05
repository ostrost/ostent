// +build linux

package types

import "fmt"

func RAMFields(ram RAM) []NameString {
	return []NameString{
		{"memory-free", fmt.Sprintf("%d", ram.Raw.Free)},
		{"memory-used", fmt.Sprintf("%d", ram.Raw.ActualUsed)},
		{"memory-buffered", fmt.Sprintf("%d", ram.Extra1)},
		{"memory-cached", fmt.Sprintf("%d", ram.Extra2)},
	}
}
