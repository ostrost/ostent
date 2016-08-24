// Package ostent is the library part of ostent cmd.
package ostent

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	metrics "github.com/rcrowley/go-metrics"

	"github.com/ostrost/ostent/internal/plugins/outputs/ostent"
	"github.com/ostrost/ostent/params"
	"github.com/ostrost/ostent/templateutil"
)

// IndexData is a data map for templates and marshalling.
// Keys (even abbrevs eg CPU) intentionally start with lowercase.
type IndexData map[string]interface{}

type renderFunc func(*params.Params) interface{}

type memoryValues struct{ Total, Free, Used uint64 }

type IndexRegistry struct{ Registry metrics.Registry }

type UpgradeInfo struct {
	RWMutex       sync.RWMutex
	LatestVersion string
}

func (ui *UpgradeInfo) Set(lv string) {
	ui.RWMutex.Lock()
	defer ui.RWMutex.Unlock()
	ui.LatestVersion = lv
}

func (ui *UpgradeInfo) Get() string {
	ui.RWMutex.RLock()
	s := ui.LatestVersion
	ui.RWMutex.RUnlock()
	if s == "" {
		return ""
	}
	return s + " release available"
}

var (
	news          map[string]renderFunc
	OstentUpgrade = new(UpgradeInfo)
	Reg1s         *IndexRegistry
)

func init() {
	reg := metrics.NewRegistry()
	Reg1s = &IndexRegistry{Registry: reg}

	// The keys with old collectors. New collectors must provide the same set.
	// olds := map[string]struct{}{
	// 	"cpu":   {},
	// 	"df":    {},
	// 	"la":    {},
	// 	"mem":   {},
	// 	"netio": {},
	//
	// 	"procs": {}, // special case
	// }
	news = map[string]renderFunc{
		"cpu": ostent.Output.CopyCPU,
		"df":  ostent.Output.CopyDisk,
		// "la" is copied with ostent.Output.CopySO
		"mem":   ostent.Output.CopyMem,
		"netio": ostent.Output.CopyNet,
		"procs": ostent.Output.CopyProc,
	}
}

func Updates(req *http.Request, para *params.Params) (IndexData, bool, error) {
	data := IndexData{}
	if req != nil {
		if err := para.Decode(req); err != nil {
			return data, false, err
		}
		// data features "params" only when req is not nil (new request).
		// So updaters do not read data for it, but expect non-nil para as an argument.
		data["params"] = para
	}

	for key, dataFunc := range news {
		data[key] = dataFunc(para)
	}
	sodup, sola := ostent.Output.CopySO(para)
	data["system_ostent"] = sodup
	data["la"] = sola
	return data, true, nil
}

type ServeSSE struct {
	logRequests bool
}

type ServeWS struct {
	ServeSSE
	logger logger
}

type ServeIndex struct {
	ServeWS
	StaticData
	IndexTemplate *templateutil.LazyTemplate
}

type StaticData struct {
	TAGGEDbin     bool
	Distrib       string
	OstentVersion string
}

func NewServeSSE(logRequests bool) ServeSSE {
	return ServeSSE{logRequests: logRequests}
}

func NewServeWS(se ServeSSE, lg logger) ServeWS { return ServeWS{ServeSSE: se, logger: lg} }

func NewServeIndex(sw ServeWS, template *templateutil.LazyTemplate, sd StaticData) ServeIndex {
	return ServeIndex{ServeWS: sw, StaticData: sd, IndexTemplate: template}
}

// Index renders index page.
func (si ServeIndex) Index(w http.ResponseWriter, r *http.Request) {
	para := params.NewParams()
	data, _, err := Updates(r, para)
	if err != nil {
		if _, ok := err.(params.RenamedConstError); ok {
			http.Redirect(w, r, err.Error(), http.StatusFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	si.IndexTemplate.Apply(w, struct {
		StaticData
		OstentUpgrade string
		Exporting     ExportingList
		Data          IndexData
	}{
		StaticData:    si.StaticData,
		OstentUpgrade: OstentUpgrade.Get(),
		Exporting:     Exporting, // from ./ws.go
		Data:          data,
	})
}

type SSE struct {
	Writer      http.ResponseWriter // points to the writer
	Params      *params.Params
	SentHeaders bool
	Errord      bool
}

// ServeHTTP is a regular serve func except the first argument,
// passed as a copy, is unused. sse.Writer is there for writes.
func (sse *SSE) ServeHTTP(_ http.ResponseWriter, r *http.Request) {
	w := sse.Writer
	data, _, err := Updates(r, sse.Params)
	if err != nil {
		if _, ok := err.(params.RenamedConstError); ok {
			http.Redirect(w, r, err.Error(), http.StatusFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	text, err := json.Marshal(data)
	if err != nil {
		sse.Errord = true
		// what would http.Error do
		if sse.SetHeader("Content-Type", "text/plain; charset=utf-8") {
			w.WriteHeader(http.StatusInternalServerError)
		}
		fmt.Fprintln(w, err.Error())
		return
	}
	sse.SetHeader("Content-Type", "text/event-stream")
	if _, err := w.Write(append(append([]byte("data: "), text...), []byte("\n\n")...)); err != nil {
		sse.Errord = true
	}
}

func (sse *SSE) SetHeader(name, value string) bool {
	if sse.SentHeaders {
		return false
	}
	sse.SentHeaders = true
	sse.Writer.Header().Set(name, value)
	return true
}

// IndexSSE serves SSE updates.
func (ss ServeSSE) IndexSSE(w http.ResponseWriter, r *http.Request) {
	sse := &SSE{Writer: w, Params: params.NewParams()}
	if LogHandler(ss.logRequests, sse).ServeHTTP(nil, r); sse.Errord {
		return
	}
	for { // loop is log-requests-free
		sleepTilNextSecond()
		if sse.ServeHTTP(nil, r); sse.Errord {
			break
		}
	}
}
