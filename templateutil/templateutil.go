// Package templateutil features LazyTemplate.
package templateutil

import (
	"bytes"
	"html/template"
	"io"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

// NewLT constructs LazyTemplate.
func NewLT(readfunc readFunc, infofunc infoFunc, filename string) *LazyTemplate {
	return &LazyTemplate{
		readFunc: readfunc,
		infoFunc: infofunc,
		filename: filename,
	}
}

type readFunc func(string) ([]byte, error)      // shortcut for Asset-type func
type infoFunc func(string) (os.FileInfo, error) // shortcut for AssetInfo-type func

// LazyTemplate encloses template.Template.
// Lazy parse
// , always clone for bin templates
// , sometimes re-parse for dev-bins
// . NewLT is the constructor.
type LazyTemplate struct {
	Mutex    sync.Mutex // protects everything
	Template *template.Template

	// arguments to NewLT (all required)
	readFunc readFunc
	infoFunc infoFunc
	filename string

	// operationals
	nonDev     bool
	devModTime time.Time
	err        error
}

// MustInit is a Must func for LazyTemplate.
func MustInit(lt *LazyTemplate) {
	lt.Mutex.Lock()
	defer lt.Mutex.Unlock()
	lt.init()
	template.Must(lt.Template, lt.err)
}

func (lt *LazyTemplate) init() { // init is internal and lock-free.
	if lt.err != nil {
		return
	}
	if lt.nonDev && lt.Template != nil {
		// allgood#1: non-dev mode & have .Template
		return
	}
	text, err := lt.readFunc(lt.filename)
	if err != nil {
		lt.err = err
		return
	}
	if info, err := lt.infoFunc(lt.filename); err != nil {
		lt.err = err
		return
	} else if modtime := info.ModTime(); modtime == time.Unix(1400000000, 0) {
		lt.nonDev = true
	} else {
		if lt.Template != nil && modtime == lt.devModTime {
			// allgood#2: dev mode + modtime did not change
			return
		}
		lt.devModTime = modtime
	}
	lt.Template, lt.err = template.
		New(lt.filename).
		Option("missingkey=error").
		Delims("[[", "]]").
		Parse(string(text))
}

// Apply executes enclosed template into w.
func (lt *LazyTemplate) Apply(w http.ResponseWriter, data interface{}) {
	clone, err := func() (*template.Template, error) {
		lt.Mutex.Lock()
		defer lt.Mutex.Unlock()
		if lt.init(); lt.err != nil {
			return nil, lt.err
		}
		return lt.Template.Clone()
	}()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	buf := new(bytes.Buffer)
	if err := clone.Execute(buf, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Content-Length", strconv.Itoa(buf.Len()))
	if _, err := io.Copy(w, buf); /* or w.Write(buf.Bytes()) */ err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
