// Package assetutil provides asset utilities.
package assetutil

import "time"

// TimeInfo is for *Asset{Info,Read}Func: a reduced os.FileInfo.
type TimeInfo interface {
	ModTime() time.Time
}
