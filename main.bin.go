// +build bin

package main

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"

	update "github.com/inconshreveable/go-update"
	"github.com/rcrowley/goagain"
	"github.com/spf13/cobra"

	"github.com/ostrost/ostent/cmd"
	"github.com/ostrost/ostent/ostent"
	"github.com/ostrost/ostent/share/templates"
)

var (
	// BootTime is the boot time.
	BootTime = time.Now()
	// AssetAltModTimeFunc returns BootTime to be asset ModTime.
	AssetAltModTimeFunc = func() time.Time { return BootTime }

	// HadUpgrade is true after an upgrade.
	HadUpgrade = new(bool)

	upLog = log.New(os.Stderr, "[ostent upgrade] ", log.LstdFlags)
	wrLog = log.New(os.Stderr, fmt.Sprintf("[%d][ostent webserver] ", os.Getpid()),
		log.LstdFlags|log.Lmicroseconds)
)

func init() {
	log.SetPrefix(fmt.Sprintf("[%d][ostent] ", os.Getpid()))
	// goagain logging is useless without pid prefix
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
}

func OstentRunE(cmd *cobra.Command, args []string) error {
	OstentUpgradeRun(false)

	if !*HadUpgrade { // unless just had an upgrade and gonna relaunch anyway
		ostent.RunBackground()
	}
	return OstentWebserverRunE()
}

func OstentWebserverServe(listen net.Listener) {
	go UntilUpgrade()
	go func() {
		templates.InitTemplates(nil) // preventive
		// sequential: Serve must wait for InitTemplates
		Serve(listen, true, nil) // true stands for taggedbin
	}()
}

func OstentWebserverRunE() error {
	listen, err := goagain.Listener()
	if err != nil {
		listen, err = net.Listen("tcp", cmd.OstentBind.String())
		if err != nil {
			return err
		}

		if *HadUpgrade { // from ./upgrade.bin.go
			go func() { // delayed kill
				time.Sleep(time.Second) // not before goagain.Wait
				GoAgain()
				// goagain.ForkExec(listen)
			}()
		} else {
			OstentWebserverServe(listen)
		}
	} else {
		OstentWebserverServe(listen)

		if err := goagain.Kill(); err != nil {
			wrLog.Fatalln(err)
		}
	}

	if _, err := goagain.Wait(listen); err != nil { // signals before won't be catched
		wrLog.Fatalln(err)
	}

	// shutting down

	if ostent.Connections.Reload() {
		time.Sleep(time.Second) // wait for an affect
	}

	if err := listen.Close(); err != nil {
		wrLog.Fatalln(err)
	}
	time.Sleep(time.Second)
	return nil
}

func main() {
	cmd.OstentCmd.RunE = OstentRunE
	cmd.Execute()
}

func NewerVersion() (string, error) {
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

func OstentUpgradeUpgrade() bool {
	newVersion, err := NewerVersion()
	if err != nil {
		upLog.Print(err)
		return false
	}
	if newVersion == "" || newVersion[0] != 'v' {
		upLog.Printf("Version unexpected: %q", newVersion)
		return false
	}
	if newVersion == "v"+ostent.VERSION {
		// upLog.Printf("Current version %q is up to date", ostent.VERSION)
		return false
	}
	upLog.Printf("Upgrade available: release %s\n", newVersion[1:])
	if cmd.DonotUpgrade {
		upLog.Print("Upgrade not applied, as requested")
		return false
	}
	upLog.Printf("Upgrading from current version %s\n", ostent.VERSION)
	url := fmt.Sprintf("https://github.com/ostrost/ostent/releases/download/%s/%s.%s",
		newVersion, strings.Title(runtime.GOOS), RuntimeMach())
	resp, err := http.Get(url)
	if err != nil {
		upLog.Print(err)
		return false
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK /* 200 */ &&
		resp.StatusCode != http.StatusNotModified /* 304 */ {
		upLog.Printf("Status not good: %d", resp.StatusCode)
		return false
	}
	/*
		if ct := resp.Header().Get("Content-Type"); ct != "application/octet-stream" {
			upLog.Printf("Content-Type is not application/octet-stream: %q", ct)
			return false
		} // */
	err = update.Apply(resp.Body, update.Options{})
	if err != nil {
		upLog.Print(err)
		return false
	}
	upLog.Printf("Upgraded successfully to release %s", newVersion[1:])
	*HadUpgrade = true
	GoAgain()
	return true
}

// GoAgain kills the current process with USR2.
func GoAgain() { syscall.Kill(os.Getpid(), syscall.SIGUSR2) }

// GoneAgain return whether the process is restarted with goagain.
func GoneAgain() bool { return os.Getenv("GOAGAIN_PPID") != "" }

func UntilUpgrade() {
	seed := time.Now().UTC().UnixNano()
	random := rand.New(rand.NewSource(seed))

	wait := time.Hour
	wait += time.Duration(random.Int63n(int64(wait))) // 1.5 +- 0.5 h
	for {
		time.Sleep(wait)
		if OstentUpgradeUpgrade() {
			break
		}
	}
}

func OstentUpgradeRun(isCommand bool) {
	if cmd.UpgradeLater {
		return
	}
	if GoneAgain() {
		// initial upgrade skipped after gone again
		return
	}
	if isCommand {
		upLog.Println("Checking for upgrades")
	} else {
		upLog.Println("Initial check for upgrades; run with --ugradelater to delay")
	}
	OstentUpgradeUpgrade()
}
