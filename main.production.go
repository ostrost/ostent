// +build production

package main

import (
	"flag"
	"net"

	"github.com/ostrost/ostent/assets"
	"github.com/ostrost/ostent/commands"
	_ "github.com/ostrost/ostent/commands/ostent"
	_ "github.com/ostrost/ostent/init-stdlogfilter"
	"github.com/ostrost/ostent/ostent"
	"github.com/ostrost/ostent/share/templates"
)

var (
	// AssetInfoFunc is for wrapping bindata's AssetInfo func.
	AssetInfoFunc = assets.ProductionAssetInfoFunc
	// AssetReadFunc is for wrapping bindata's Asset func.
	AssetReadFunc = assets.ProductionAssetReadFunc
)

func init() {
	commands.InitStdLog()
}

func main() {
	flag.Usage = commands.UsageFunc(flag.CommandLine)
	webserver := commands.NewWebserver(8050).AddCommandLine()
	upgrade := commands.NewUpgrade().AddCommandLine()
	flag.Parse()

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
			Serve(listen, true, nil) // true stands for production
		}()
	}
	upgrade.FirstUpgradeStopper = webserver.GoneAgain // initial upgrade skipped after gone again
	upgrade.AfterUpgradeFunc = webserver.GoAgain
	upgrade.Run()

	webserver.FirstRunFunc = upgrade.HadUpgrade
	if !upgrade.HadUpgrade() {
		// RunBackground unless just had an upgrade and gonna relaunch anyway
		ostent.RunBackground(periodFlag)
	}
	webserver.Run()
}
