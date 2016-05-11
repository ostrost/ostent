package main

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/http/pprof"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/facebookgo/grace/gracehttp"
	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
	"github.com/spf13/cobra"

	"github.com/blang/semver" // alt semver: "github.com/Masterminds/semver"
	"github.com/ostrost/ostent/cmd"
	"github.com/ostrost/ostent/ostent"
	"github.com/ostrost/ostent/share/assets"
	"github.com/ostrost/ostent/share/templates"
)

func main() {
	cmd.OstentCmd.RunE = func(*cobra.Command, []string) error {
		return Serve(cmd.OstentBind.String())
	}
	cmd.Execute()
}

var (
	// NoUpgradeCheck is the flag value.
	NoUpgradeCheck bool

	logRequests bool // flag
	taggedBin   bool // whether the build features bin tag
)

func init() {
	taggedBin = assets.AssetAltModTimeFunc != nil

	cmd.OstentCmd.Flags().BoolVar(&logRequests, "log-requests", !taggedBin,
		"Whether to log webserver requests")
	cmd.OstentCmd.Flags().BoolVar(&NoUpgradeCheck, "noupgradecheck", false,
		"Off periodic upgrade check")
	ostent.AddBackground(ostent.CollectLoop)

	var err error
	CurrentV, err = NewSemver(ostent.VERSION)
	if err != nil { // CurrentV stay nil
		log.Printf("Current semver parse error: %s\n", err)
	}
}

// CurrentV is a *semver.Version from ostent.VERSION.
var CurrentV *semver.Version

func NewSemver(s string) (*semver.Version, error) {
	v, err := semver.New(s)
	if err == nil {
		return v, nil
	}
	if err.Error() == "No Major.Minor.Patch elements found" &&
		len(strings.SplitN(s, ".", 3)) == 2 {
		return semver.New(s + ".0")
	}
	return nil, err
}

// Serve constructs a *http.Server to (gracefully) Serve. Routes are set here.
func Serve(laddr string) error {
	if !NoUpgradeCheck && CurrentV != nil {
		go UntilUpgradeCheck(CurrentV)
	}
	ostent.RunBackground()
	templates.InitTemplates()

	r, achain := ostent.NewServery(taggedBin)

	ostentLog := log.New(os.Stderr, "[ostent] ", 0)
	errlog, errclose := ostent.NewErrorLog()
	defer errclose()

	var (
		serve1 = ostent.NewServeSSE(logRequests, cmd.DelayFlags)
		serve2 = ostent.NewServeWS(serve1, errlog)
		serve3 = ostent.NewServeIndex(serve2, taggedBin, templates.IndexTemplate)
		serve4 = ostent.ServeAssets{
			Log:                 ostentLog,
			AssetFunc:           assets.Asset,
			AssetInfoFunc:       assets.AssetInfo,
			AssetAltModTimeFunc: assets.AssetAltModTimeFunc,
		}
	)

	achainT := ostent.HandleThen(achain.Then)
	paramsT := ostent.ParamsFunc(achain.Append(context.ClearHandler).Then)
	var (
		panicp = "/panic"
		panics = func(http.ResponseWriter, *http.Request) { panic(panicp) }
		panicr = [2]string{panicp, "GET HEAD"}
		pprofr = [2]string{"/debug/pprof/:name", "GET HEAD POST"}
	)

	routes := map[[2]string]httprouter.Handle{
		{"/index.sse", "GET"}: ostent.HandleFunc(serve1.IndexSSE),
		{"/index.ws", "GET"}:  ostent.HandleFunc(serve2.IndexWS),
		{"/", "GET HEAD"}:     achainT(serve3.Index),
		panicr:                achainT(panics),
		pprofr:                paramsT(pprofHandle),
	}
	for _, path := range assets.AssetNames() {
		p := "/" + path
		if path != "favicon.ico" && path != "robots.txt" {
			p = "/" + ostent.VERSION + "/" + path // the Version prefix
		}
		routes[[2]string{p, "GET HEAD"}] = ostent.HandleThen(achain.Append(
			context.ClearHandler,
			ostent.AddAssetPathContextFunc(path),
		).Then)(serve4.Serve)
	}
	if taggedBin { // no panicr or pprofr in bin
		delete(routes, panicr)
		delete(routes, pprofr)
	}

	for x, handle := range routes {
		for _, m := range strings.Split(x[1], " ") {
			r.Handle(m, x[0], handle)
		}
	}
	return gracehttp.Serve(&http.Server{
		Addr:     laddr,
		ErrorLog: errlog,
		Handler:  ostent.LogHandler(logRequests, r),
	})
}

// UntilUpgradeCheck waits for an upgrade and returns.
func UntilUpgradeCheck(cv *semver.Version) {
	if UpgradeCheck(cv) {
		return
	}

	seed := time.Now().UTC().UnixNano()
	random := rand.New(rand.NewSource(seed))

	wait := time.Hour
	wait += time.Duration(random.Int63n(int64(wait))) // 1.5 +- 0.5 h
	for {
		time.Sleep(wait)
		if UpgradeCheck(cv) {
			break
		}
	}
}

// UpgradeCheck does upgrade check and returns true if an upgrade is available.
func UpgradeCheck(cv *semver.Version) bool {
	newVersion, err := NewerVersion()
	if err != nil {
		log.Printf("Upgrade check error: %s\n", err)
		return false
	}
	if newVersion == "" {
		log.Printf("Upgrade check: version is empty\n")
		return false
	}
	nv, err := NewSemver(newVersion)
	if err != nil {
		log.Printf("Semver parse error: %s\n", err)
		return false
	}
	if !nv.GT(*cv) {
		return false
	}
	log.Printf("Upgrade check: %s release available\n", newVersion)
	ostent.OstentUpgrade.Set(newVersion)
	return true
}

// NewerVersion checks GitHub for the latest ostent version.
// Return is in form of "\d.*" (sans "^v").
func NewerVersion() (string, error) {
	// 1. https://github.com/ostrost/ostent/releases/latest // redirects, NOT followed
	// 2. https://github.com/ostrost/ostent/releases/v...   // Redirect location
	// 3. return "v..." // basename of the location

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
	v := filepath.Base(redir.url.Path)
	if len(v) == 0 || v[0] != 'v' {
		return "", fmt.Errorf("Unexpected version from GitHub: %q", v)
	}
	return v[1:], nil
}

func pprofHandle(w http.ResponseWriter, r *http.Request) {
	params, err := ostent.ContextParams(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	name := params.ByName("name")
	var handler func(http.ResponseWriter, *http.Request)
	switch name {
	case "cmdline":
		handler = pprof.Cmdline
	case "profile":
		handler = pprof.Profile
	case "symbol":
		handler = pprof.Symbol
	case "trace":
		handler = pprof.Trace
	default:
		handler = pprof.Index
	}
	if handler == nil {
		http.NotFound(w, r)
		return
	}
	handler(w, r)
}
