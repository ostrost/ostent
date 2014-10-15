// +build production

package main

import (
	"flag"
	"net"

	"github.com/ostrost/ostent"
	"github.com/ostrost/ostent/commands"
	"github.com/ostrost/ostent/share/templates"
)

func init() {
	commands.InitStdLog()
}

func main() {
	flag.Usage = commands.UsageFunc(flag.CommandLine)
	webserver := commands.FlagSetNewWebserver(flag.CommandLine)
	// version := commands.FlagSetNewVersion(flag.CommandLine)
	upgrade := commands.FlagSetNewUpgrade(flag.CommandLine)
	flag.Parse()
	defer commands.Defaults()()

	if errd := commands.ArgCommands(); errd { // explicit commands
		return
	}
	// if version.Flag { version.Run(); return }

	webserver.ShutdownFunc = ostent.Connections.Reload
	webserver.ServeFunc = func(listen net.Listener) {
		go upgrade.UntilUpgrade()
		go templates.InitTemplates() // preventive
		go Serve(listen, true, nil)  // true stands for production
	}
	upgrade.FirstUpgradeStopper = webserver.GoneAgain // initial upgrade skipped after gone again
	upgrade.AfterUpgradeFunc = webserver.GoAgain
	upgrade.Run()

	webserver.FirstRunFunc = upgrade.HadUpgrade
	if !upgrade.HadUpgrade() {
		// start the background job(s) unless just had an upgrade and gonna relaunch anyway
		go ostent.Loop()
		// go ostent.CollectdLoop()
	}
	webserver.Run()
}
