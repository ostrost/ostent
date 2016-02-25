// +build bin

package main

import (
	"time"

	"github.com/spf13/cobra"

	"github.com/ostrost/ostent/cmd"
)

var (
	// BootTime is the boot time.
	BootTime = time.Now()
	// AssetAltModTimeFunc returns BootTime to be asset ModTime.
	AssetAltModTimeFunc = func() time.Time { return BootTime }
)

func main() {
	cmd.OstentCmd.RunE = func(*cobra.Command, []string) error {
		return Serve(cmd.OstentBind.String(), true, nil)
	}
	cmd.Execute()
}
