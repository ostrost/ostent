// +build production

package assets
import (
	"time"
)

var STARTIME = time.Now()
func ModTime(string, string) (time.Time, error) {
	return STARTIME, nil
}
