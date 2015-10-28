// Package templateutil features LazyTemplate and TemplateWriter.
package templateutil

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strconv"
	"sync"
)

// NewLT constructs LazyTemplate.
func NewLT(readfunc Readfunc, filename string, funcmap template.FuncMap) *LazyTemplate {
	return &LazyTemplate{Readfunc: readfunc, Filename: filename, Funcmap: funcmap}
}

// Readfunc type is shortcut for Asset-type func.
type Readfunc func(string) ([]byte, error)

// LazyTemplate has a template.Template. Lazy parse, always clone.
type LazyTemplate struct {
	MU sync.Mutex // protects everything

	// arguments to NewLT (all required)
	Readfunc Readfunc
	Filename string
	Funcmap  template.FuncMap

	Template *template.Template
	Err      error
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
	if lt.Template != nil {
		return
	}
	if lt.Readfunc == nil {
		lt.Err = fmt.Errorf("templateutil: readfunc is nil for %q reading", lt.Filename)
		return
	}
	text, err := lt.Readfunc(lt.Filename)
	if err != nil {
		lt.Err = err
		return
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
	io.Copy(w, buf) // or w.Write(buf.Bytes())
}
