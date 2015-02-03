// +build !production

package assets

var rootDir string

// RootDir sets the prefix for asset files
func RootDir(dir string) { rootDir = dir }
