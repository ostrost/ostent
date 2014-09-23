// +build production

package commands

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	update "github.com/inconshreveable/go-update"
	"github.com/ostrost/ostent"
)

func (u upgrade) newerVersion() (string, error) {
	// 1. https://github.com/ostrost/ostent/releases/latest # redirects, NOT followed
	// 2. https://github.com/ostrost/ostent/releases/vX.Y.Z # Redirect location
	// 3. return "vX.Y.Z" # basename of the location

	type redirected struct {
		error
		url url.URL
	}
	checkRedirect := func(req *http.Request, _via []*http.Request) error {
		return redirected{url: *req.URL}
	}

	client := &http.Client{CheckRedirect: checkRedirect}
	resp, err := client.Get("https://github.com/ostrost/ostent/releases/latest")
	if err == nil {
		resp.Body.Close()
		return "", errors.New("The GitHub /latest page did not return a redirect.")
	}
	urlerr, ok := err.(*url.Error)
	if !ok {
		return "", err
	}
	if resp != nil && resp.Body != nil {
		resp.Body.Close()
	}
	redir, ok := urlerr.Err.(redirected)
	if !ok {
		return "", urlerr
	}
	return filepath.Base(redir.url.Path), nil
}

func (_ upgrade) mach() string {
	m := runtime.GOARCH
	if m == "amd64" {
		return "x86_64"
	} else if m == "386" {
		return "i686"
	}
	return m
}

func (up *upgrade) Upgrade() bool {
	newVersion, err := up.newerVersion()
	if err != nil {
		up.logger.Print(err)
		return false
	}
	if newVersion == "" || newVersion[0] != 'v' {
		up.logger.Printf("Version unexpected: %q", newVersion)
		return false
	}
	if newVersion == "v"+ostent.VERSION {
		up.logger.Printf("Current version %q is up to date", ostent.VERSION)
		return false
	}
	up.logger.Printf("Upgrade available: release %s\n", newVersion[1:])
	up.logger.Printf("Upgrading from current version %s\n", ostent.VERSION)
	if up.DonotUpgrade {
		up.logger.Print("Upgrade not applied, as requested")
		return false
	}
	url := fmt.Sprintf("https://github.com/ostrost/ostent/releases/download/%s/%s.%s", newVersion, strings.Title(runtime.GOOS), up.mach())
	// url = fmt.Sprintf("http://127.0.0.1:8000/%s.%s", strings.Title(runtime.GOOS), up.mach()) // testing
	err, errecov := update.New().FromUrl(url)
	if err != nil {
		up.logger.Print(err)
		return false
	}
	if errecov != nil {
		up.logger.Print(errecov)
		return false
	}
	up.logger.Printf("Upgraded successfully to release %s", newVersion[1:])
	up.hadUpgrade = true
	if up.AfterUpgradeFunc != nil {
		up.AfterUpgradeFunc()
	}
	return true
}

func (up *upgrade) UntilUpgrade() {
	wait := time.Hour
	wait += time.Duration(rand.Int63n(int64(wait))) // 1.5 +- 0.5 h
	// wait = time.Second * 20 // testing
	for {
		select {
		case <-time.After(wait):
			if up.Upgrade() { // (true)
				break
			}
		}
	}
}

func (up upgrade) Run() {
	if up.UpgradeLater {
		return
	}
	if up.FirstUpgradeStopper != nil && up.FirstUpgradeStopper() {
		return
	}
	if up.isCommand {
		up.logger.Println("Checking for upgrades")
	} else {
		up.logger.Println("Initial check for upgrades; run with -ugradelater to disable")
	}
	up.Upgrade()
}

type upgrade struct {
	DonotUpgrade        bool
	UpgradeLater        bool
	logger              *loggerWriter
	hadUpgrade          bool
	FirstUpgradeStopper func() bool
	AfterUpgradeFunc    func()
	isCommand           bool
}

func (up upgrade) HadUpgrade() bool {
	return up.hadUpgrade
}

func FlagSetNewUpgrade(fs *flag.FlagSet) *upgrade { // fs better be flag.CommandLine
	up := upgrade{logger: &loggerWriter{log.New(os.Stderr, "[ostent upgrade] ", log.LstdFlags)}}
	fs.BoolVar(&up.UpgradeLater, "upgradelater", false, "Upgrade later.")
	return &up
}

func upgradeCommand(fs *flag.FlagSet, arguments []string) (sub, error, []string) {
	up := FlagSetNewUpgrade(fs)
	up.isCommand = true
	// this Var bound to fs, thus doesn't make it into flag.CommandLine set
	fs.BoolVar(&up.DonotUpgrade, "n", false, "Do not upgrade, just log if there's an upgrade.")
	fs.SetOutput(up.logger)
	err := fs.Parse(arguments)
	return up.Run, err, fs.Args()
}

func init() {
	AddCommand("upgrade", upgradeCommand)
}
