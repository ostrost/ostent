// +build !production

package assets
import (
	"os"
	"time"
	"path/filepath"
)

var stat_fails = false
// TODO mutex this

func ModTime(prefix, path string) (time.Time, error) {
	now := time.Now()
	if stat_fails {
		return now, nil
	}
	fi, err := os.Stat(filepath.Join(prefix, path))
	if err != nil {
		stat_fails = true
		return now, err
	}
	return fi.ModTime(), nil
}
