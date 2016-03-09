package main

import (
	"errors"
	"log"
	"math/rand"
	"net/http"
	"net/http/pprof"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/facebookgo/grace/gracehttp"
	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"github.com/spf13/cobra"

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

// NoUpgradeCheck is the flag value.
var NoUpgradeCheck bool

func init() {
	cmd.OstentCmd.Flags().BoolVar(&NoUpgradeCheck, "noupgradecheck", false,
		"Off periodic upgrade check")
	ostent.AddBackground(ostent.CollectLoop)
}

// Serve constructs a *http.Server to (gracefully) Serve. Routes are set here.
func Serve(laddr string) error {
	if !NoUpgradeCheck {
		go UntilUpgradeCheck()
	}
	ostent.RunBackground()
	templates.InitTemplates()

	taggedbin := assets.AssetAltModTimeFunc != nil
	r, achain, access := ostent.NewServery(taggedbin)

	ostentLog := log.New(os.Stderr, "[ostent] ", 0)
	errlog, errclose := ostent.NewErrorLog()
	defer errclose()

	var (
		serve1 = ostent.NewServeSSE(access, cmd.DelayFlags)
		serve2 = ostent.NewServeWS(serve1, errlog)
		serve3 = ostent.NewServeIndex(serve2, taggedbin, templates.IndexTemplate)
		serve4 = ostent.ServeAssets{
			Log:                 ostentLog,
			AssetFunc:           assets.Asset,
			AssetInfoFunc:       assets.AssetInfo,
			AssetAltModTimeFunc: assets.AssetAltModTimeFunc,
		}
	)

	m, n, chain0 := ostent.NewRoute, ostent.NewHandle, alice.New()

	var (
		panicp = "/panic"
		panics = func(http.ResponseWriter, *http.Request) { panic(panicp) }
		panicr = m(panicp, r.GET, r.HEAD)
	)

	routes := map[*ostent.Route]ostent.Handle{
		m("/index.sse", r.GET): n(chain0, serve1.IndexSSE),
		m("/index.ws", r.GET):  n(chain0, serve2.IndexWS),
		m("/", r.GET, r.HEAD):  n(achain, serve3.Index),
		panicr:                 n(achain, panics),
	}
	for _, path := range assets.AssetNames() {
		p := "/" + path
		if path != "favicon.ico" && path != "robots.txt" {
			p = "/" + ostent.VERSION + "/" + path // the Version prefix
		}
		rr := m(p, r.GET, r.HEAD)
		rr.Asset = true
		routes[rr] = n(achain.Append(context.ClearHandler,
			ostent.AddAssetPathContextFunc(path)), serve4.Serve)
	}
	if taggedbin { // no panicp in bin
		delete(routes, panicr)
	}

	// now bind
	ostent.ApplyRoutes(r, routes, nil)
	if !taggedbin {
		handle := PprofServe{achain}.Serve
		r.GET("/debug/pprof/:name", handle)
		r.HEAD("/debug/pprof/:name", handle)
		r.POST("/debug/pprof/:name", handle)
	}

	return gracehttp.Serve(&http.Server{
		Addr:     laddr,
		ErrorLog: errlog,
		Handler:  r,
	})
}

// UntilUpgradeCheck waits for an upgrade and returns.
func UntilUpgradeCheck() {
	if UpgradeCheck() {
		return
	}

	seed := time.Now().UTC().UnixNano()
	random := rand.New(rand.NewSource(seed))

	wait := time.Hour
	wait += time.Duration(random.Int63n(int64(wait))) // 1.5 +- 0.5 h
	for {
		time.Sleep(wait)
		if UpgradeCheck() {
			break
		}
	}
}

// UpgradeCheck does upgrade check and returns true if an upgrade is available.
func UpgradeCheck() bool {
	newVersion, err := NewerVersion()
	if err != nil {
		log.Printf("Upgrade check error: %s\n", err)
		return false
	}
	if newVersion == "" || newVersion[0] != 'v' {
		log.Printf("Upgrade check error: version unexpected: %q\n", newVersion)
		return false
	}
	if newVersion == "v"+ostent.VERSION {
		return false
	}
	log.Printf("Upgrade check: %s release available\n", newVersion[1:])
	ostent.OstentUpgrade.Set(newVersion[1:])
	return true
}

// NewerVersion checks GitHub for the latest ostent version in form of "v...".
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
	return filepath.Base(redir.url.Path), nil
}

// Serve is PprofServe handler.
func (ps PprofServe) Serve(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
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
	ps.Chain.ThenFunc(handler).ServeHTTP(w, r)
}

// PprofServe is pprof data serving handler.
type PprofServe struct{ Chain alice.Chain }
