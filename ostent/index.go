// Package ostent is the library part of ostent cmd.
package ostent

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/ostrost/ostent/internal/plugins/outputs/ostent"
	"github.com/ostrost/ostent/params"
	"github.com/ostrost/ostent/templateutil"
)

// IndexData is a data map for templates and marshalling.
// Keys (even abbrevs eg CPU) intentionally start with lowercase.
type IndexData map[string]interface{}

type memoryValues struct{ Total, Free, Used uint64 }

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

var OstentUpgrade = new(UpgradeInfo)

func Updates(req *http.Request, para *params.Params) (IndexData, bool, error) {
	data := IndexData{}
	if decoded, ok := req.Context().Value(crequestDecoded).(bool); !ok || !decoded {
		if err := para.Decode(req); err != nil {
			return data, false, err
		}
		// data features "params" only when req is not nil (new request).
		// So updaters do not read data for it, but expect non-nil para as an argument.
		data["params"] = para
	}

	up, _ := req.Context().Value(coutputUpdate).(*ostent.Update)
	if up == nil { // http request (not ws)
		up = lastCopy.get()
	}
	if up == nil { // before first collection
		// "system_ostent" is expected to be an object (not an array)
		data["system_ostent"] = map[string]string{}
		return data, true, nil
	}

	data["system_ostent"] = ostent.Output.CopySO(up, para)
	for key, dataFunc := range map[string]func(*ostent.Update, *params.Params) interface{}{
		"cpu":   ostent.Output.CopyCPU,
		"df":    ostent.Output.CopyDisk,
		"la":    ostent.Output.CopyLA,
		"mem":   ostent.Output.CopyMem,
		"netio": ostent.Output.CopyNet,
		"procs": ostent.Output.CopyProc,
	} {
		type list struct{ List interface{} }
		data[key] = list{dataFunc(up, para)}
	}
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

	data["Distrib"] = si.StaticData.Distrib
	data["Exporting"] = exportingCopy() // from ./ws.go
	data["OstentUpgrade"] = OstentUpgrade.Get()
	data["OstentVersion"] = si.StaticData.OstentVersion
	data["TAGGEDbin"] = si.StaticData.TAGGEDbin

	si.IndexTemplate.Apply(w, struct{ Data IndexData }{Data: data})
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

		now := time.Now()
		time.Sleep(now.Truncate(time.Second).Add(time.Second).Sub(now))
		// TODO redo with some channel to receive from, pushed in lastCopy by CollectLoop or something

		if sse.ServeHTTP(nil, r); sse.Errord {
			break
		}
	}
}
