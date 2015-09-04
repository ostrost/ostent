// +build bin

package commands

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	update "github.com/inconshreveable/go-update"

	"github.com/ostrost/ostent/commands/extpoints"
	"github.com/ostrost/ostent/ostent"
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

func RuntimeMach() string {
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
		up.Log.Print(err)
		return false
	}
	if newVersion == "" || newVersion[0] != 'v' {
		up.Log.Printf("Version unexpected: %q", newVersion)
		return false
	}
	if newVersion == "v"+ostent.VERSION {
		// up.Log.Printf("Current version %q is up to date", ostent.VERSION)
		return false
	}
	up.Log.Printf("Upgrade available: release %s\n", newVersion[1:])
	up.Log.Printf("Upgrading from current version %s\n", ostent.VERSION)
	if up.DonotUpgrade {
		up.Log.Print("Upgrade not applied, as requested")
		return false
	}
	url := fmt.Sprintf("https://github.com/ostrost/ostent/releases/download/%s/%s.%s",
		newVersion, strings.Title(runtime.GOOS), RuntimeMach())
	resp, err := http.Get(url)
	if err != nil {
		up.Log.Print(err)
		return false
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK /* 200 */ &&
		resp.StatusCode != http.StatusNotModified /* 304 */ {
		up.Log.Printf("Status not good: %d", resp.StatusCode)
		return false
	}
	/*
		if ct := resp.Header().Get("Content-Type"); ct != "application/octet-stream" {
			up.Log.Printf("Content-Type is not application/octet-stream: %q", ct)
			return false
		} // */
	err = update.Apply(resp.Body, update.Options{})
	if err != nil {
		up.Log.Print(err)
		return false
	}
	up.Log.Printf("Upgraded successfully to release %s", newVersion[1:])
	up.hadUpgrade = true
	if up.AfterUpgradeFunc != nil {
		up.AfterUpgradeFunc()
	}
	return true
}

func (up *upgrade) UntilUpgrade() {
	seed := time.Now().UTC().UnixNano()
	random := rand.New(rand.NewSource(seed))

	wait := time.Hour
	wait += time.Duration(random.Int63n(int64(wait))) // 1.5 +- 0.5 h
	for {
		time.Sleep(wait)
		if up.Upgrade() {
			break
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
		up.Log.Println("Checking for upgrades")
	} else {
		up.Log.Println("Initial check for upgrades; run with -ugradelater to disable")
	}
	up.Upgrade()
}

type upgrade struct {
	DonotUpgrade        bool
	UpgradeLater        bool
	Log                 *extpoints.Log
	hadUpgrade          bool
	FirstUpgradeStopper func() bool
	AfterUpgradeFunc    func()
	isCommand           bool
}

func (up upgrade) HadUpgrade() bool {
	return up.hadUpgrade
}

func NewUpgrade(loggerOptions ...extpoints.SetupLog) *upgrade {
	return &upgrade{
		Log: NewLog("[ostent upgrade] ", loggerOptions...),
	}
}

func (u Upgrades) SetupCommand(fs *flag.FlagSet, loggerOptions ...extpoints.SetupLog) (extpoints.CommandHandler, io.Writer) {
	up := NewUpgrade(loggerOptions...)
	up.isCommand = true
	fs.BoolVar(&up.DonotUpgrade, "n", false, "Do not upgrade, just log if there's an upgrade.")
	return up.Run, up.Log
}

func (up *upgrade) SetupFlagSet(cli *flag.FlagSet) extpoints.CommandLineHandler {
	cli.BoolVar(&up.UpgradeLater, "upgradelater", false, "Upgrade later.")
	return nil
}

func (up *upgrade) AddCommandLine() *upgrade {
	extpoints.CommandLines.Register(up, "upgrade")
	return up
}

type Upgrades struct{}

func init() {
	extpoints.Commands.Register(Upgrades{}, "upgrade")
}
