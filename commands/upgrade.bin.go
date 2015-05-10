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
		up.logger.Println("Checking for upgrades")
	} else {
		up.logger.Println("Initial check for upgrades; run with -ugradelater to disable")
	}
	up.Upgrade()
}

type upgrade struct {
	DonotUpgrade        bool
	UpgradeLater        bool
	logger              *Logger
	hadUpgrade          bool
	FirstUpgradeStopper func() bool
	AfterUpgradeFunc    func()
	isCommand           bool
}

func (up upgrade) HadUpgrade() bool {
	return up.hadUpgrade
}

func (up *upgrade) AddCommandLine() *upgrade {
	AddCommandLine(func(cli *flag.FlagSet) CommandLineHandler {
		cli.BoolVar(&up.UpgradeLater, "upgradelater", false, "Upgrade later.")
		return nil
	})
	return up
}

func NewUpgrade(loggerOptions ...SetupLogger) *upgrade {
	return &upgrade{
		logger: NewLogger("[ostent upgrade] ", loggerOptions...),
	}
}

func upgradeCommand(fs *flag.FlagSet, loggerOptions ...SetupLogger) (CommandHandler, io.Writer) {
	up := NewUpgrade(loggerOptions...)
	up.isCommand = true
	fs.BoolVar(&up.DonotUpgrade, "n", false, "Do not upgrade, just log if there's an upgrade.")
	return up.Run, up.logger
}

func init() {
	AddCommand("upgrade", upgradeCommand)
}
