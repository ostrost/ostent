// +build bin

package assets

import "time"

var (
	// BootTime is the boot time.
	BootTime = time.Now()
	// AssetAltModTimeFunc returns BootTime to be asset ModTime.
	AssetAltModTimeFunc = func() time.Time { return BootTime }
)
