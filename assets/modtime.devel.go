// +build !production

package assets

import "time"

// ModTime returns t.
func ModTime(t time.Time) time.Time { return t }

// UncompressedAssetFunc returns readFunc.
func UncompressedAssetFunc(readFunc func(string) ([]byte, error)) func(string) ([]byte, error) {
	return readFunc
}
