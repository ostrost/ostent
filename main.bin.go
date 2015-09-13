// +build bin

package main

import (
	"flag"
	"net"
	"os"
	"time"

	"github.com/ostrost/ostent/commands"
	_ "github.com/ostrost/ostent/commands/ostent"
	"github.com/ostrost/ostent/ostent"
	"github.com/ostrost/ostent/share/templates"
)

var (
	// BootTime is the boot time.
	BootTime = time.Now()
	// AssetAltModTimeFunc returns BootTime to be asset ModTime.
	AssetAltModTimeFunc = func() time.Time { return BootTime }
)

func init() {
	commands.InitStdLog()
}

func main() {
	var (
		webserver = commands.NewWebserver(8050).AddCommandLine()
		upgrade   = commands.NewUpgrade().AddCommandLine()
	)
	commands.Parse(flag.CommandLine, os.Args[1:])

	errd, atexit := commands.ArgCommands()
	defer atexit()

	if errd {
		return
	}

	webserver.ShutdownFunc = ostent.Connections.Reload
	webserver.ServeFunc = func(listen net.Listener) {
		go upgrade.UntilUpgrade()
		go func() {
			templates.InitTemplates(nil) // preventive
			// sequential: Serve must wait for InitTemplates
			Serve(listen, true, nil) // true stands for taggedbin
		}()
	}
	upgrade.FirstUpgradeStopper = webserver.GoneAgain // initial upgrade skipped after gone again
	upgrade.AfterUpgradeFunc = webserver.GoAgain
	upgrade.Run()

	webserver.FirstRunFunc = upgrade.HadUpgrade
	if !upgrade.HadUpgrade() {
		// RunBackground unless just had an upgrade and gonna relaunch anyway
		ostent.RunBackground(MinDelayFlag)
	}
	webserver.Run()
}
