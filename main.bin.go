// +build bin

package main

import (
	"flag"
	"net"
	"os"
	"sync"
	"time"

	"github.com/ostrost/ostent/commands"
	_ "github.com/ostrost/ostent/commands/ostent"
	_ "github.com/ostrost/ostent/init-stdlogfilter"
	"github.com/ostrost/ostent/ostent"
	"github.com/ostrost/ostent/share/templates"
)

var (
	// AssetInfoFunc is for wrapping bindata's AssetInfo func.
	AssetInfoFunc = BinAssetInfoFunc
	// AssetReadFunc is for wrapping bindata's Asset func.
	AssetReadFunc = BinAssetReadFunc
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
			Serve(listen, true, nil) // true stands for taggedbin
		}()
	}
	upgrade.FirstUpgradeStopper = webserver.GoneAgain // initial upgrade skipped after gone again
	upgrade.AfterUpgradeFunc = webserver.GoAgain
	upgrade.Run()

	webserver.FirstRunFunc = upgrade.HadUpgrade
	if !upgrade.HadUpgrade() {
		// RunBackground unless just had an upgrade and gonna relaunch anyway
		ostent.RunBackground(PeriodFlag)
	}
	webserver.Run()
}

var (
	// BootTime is the boot time.
	BootTime = time.Now()
	// ReadCache is a cache in AssetReadFunc.
	ReadCache struct {
		MU     sync.Mutex
		Byname map[string][]byte
	}
)

// BinAssetInfoFunc wraps bindata's AssetInfo func. ModTime is always BootTime.
func BinAssetInfoFunc(infofunc func(string) (os.FileInfo, error)) func(string) (ostent.TimeInfo, error) {
	return func(name string) (ostent.TimeInfo, error) {
		_, err := infofunc(name)
		if err != nil {
			return nil, err
		}
		return BootInfo{}, nil
	}
}

// BinAssetReadFunc wraps bindata's Asset func. Result is from cache or cached.
func BinAssetReadFunc(readfunc func(string) ([]byte, error)) func(string) ([]byte, error) {
	return func(name string) ([]byte, error) {
		return Read(readfunc, name)
	}
}

// BootInfo is a ostent.TimeInfo implementation.
type BootInfo struct{}

// ModTime returns BootTime.
func (bi BootInfo) ModTime() time.Time { return BootTime } // bi is unused

// Read returns cached readfunc result.
func Read(readfunc func(string) ([]byte, error), name string) ([]byte, error) {
	ReadCache.MU.Lock()
	defer ReadCache.MU.Unlock()
	if text, ok := ReadCache.Byname[name]; ok {
		return text, nil
	}
	text, err := readfunc(name)
	if err != nil {
		ReadCache.Byname[name] = text
	}
	return text, err
}
