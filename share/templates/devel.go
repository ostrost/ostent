// +build !production

package templates

var rootDir string

// RootDir sets the prefix for template files
func RootDir(dir string) { rootDir = dir }
