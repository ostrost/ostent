// Package templateutil features LazyTemplate and TemplateWriter.
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
func NewLT(readfunc ReadFunc, infofunc InfoFunc, filename string, funcmap template.FuncMap) *LazyTemplate {
	return &LazyTemplate{
		ReadFunc: readfunc,
		InfoFunc: infofunc,
		Filename: filename,
		Funcmap:  funcmap,
	}
}

// ReadFunc type is shortcut for Asset-type func.
type ReadFunc func(string) ([]byte, error)

// InfoFunc type is shortcut for AssetInfo-type func.
type InfoFunc func(string) (os.FileInfo, error)

// LazyTemplate has a template.Template.
// Lazy parse
// , always clone for bin templates
// , sometimes re-parse for dev-bins
// . NewLT is the constructor.
type LazyTemplate struct {
	MU sync.Mutex // protects everything

	// arguments to NewLT (all required)
	ReadFunc ReadFunc
	InfoFunc InfoFunc
	Filename string
	Funcmap  template.FuncMap

	// operationals
	NonDev     bool
	DevModTime time.Time
	Template   *template.Template
	Err        error
}

// MustInit is a Must func for LazyTemplate.
func MustInit(lt *LazyTemplate) {
	lt.MU.Lock()
	defer lt.MU.Unlock()
	lt.Init()
	template.Must(lt.Template, lt.Err)
}

// Init is internal and lock-free.
func (lt *LazyTemplate) Init() {
	if lt.Err != nil {
		return
	}
	if lt.NonDev && lt.Template != nil {
		// allgood#1: non-dev mode & have .Template
		return
	}
	text, err := lt.ReadFunc(lt.Filename)
	if err != nil {
		lt.Err = err
		return
	}
	if info, err := lt.InfoFunc(lt.Filename); err != nil {
		lt.Err = err
		return
	} else if modtime := info.ModTime(); modtime == time.Unix(1400000000, 0) {
		lt.NonDev = true
	} else {
		if lt.Template != nil && modtime == lt.DevModTime {
			// allgood#2: dev mode + modtime did not change
			return
		}
		lt.DevModTime = modtime
	}
	t := template.New(lt.Filename)
	t = t.Option("missingkey=error")
	if lt.Funcmap != nil {
		t.Funcs(lt.Funcmap)
	}
	lt.Template, lt.Err = t.Parse(string(text))
}

// Apply clones .Template to execute it into w.
func (lt *LazyTemplate) Apply(w http.ResponseWriter, data interface{}) {
	clone, err := func() (*template.Template, error) {
		lt.MU.Lock()
		defer lt.MU.Unlock()
		if lt.Init(); lt.Err != nil {
			return nil, lt.Err
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
