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

	// these accessed in InitClean thus also protected by MU
	Readfunc Readfunc         // required (argument to NewLT)
	Filename string           // required (argument to NewLT)
	Funcmap  template.FuncMap // required (argument to NewLT)

	Template *template.Template
	Err      error
}

// MustInit is a Must func for LazyTemplate.
func MustInit(l *LazyTemplate) {
	l.MU.Lock()
	defer l.MU.Unlock()
	err := l.Init()
	template.Must(l.Template, err)
}

// Init is internal. Keeps .Err and/may call InitClean.
func (l *LazyTemplate) Init() error {
	if l.Err != nil {
		return l.Err
	}
	err := l.InitClean()
	if err != nil {
		l.Err = err
	}
	return err
}

// InitClean is internal. The init without .Err handling.
func (l *LazyTemplate) InitClean() error {
	if l.Template != nil {
		return nil
	}
	if l.Readfunc == nil {
		return fmt.Errorf("templateutil: readfunc is nil for %q reading", l.Filename)
	}
	text, err := l.Readfunc(l.Filename)
	if err != nil {
		return err
	}
	t := template.New(l.Filename)
	if l.Funcmap != nil {
		t.Funcs(l.Funcmap)
	}
	t, err = t.Parse(string(text))
	if err != nil {
		return err
	}
	l.Template = t
	return nil
}

// LookupApply wraps ApplyTemplate with l.Template.Lookup(name).
func (l *LazyTemplate) LookupApply(name string, data interface{}) (*bytes.Buffer, error) {
	return l.ApplyTemplate(func() (*template.Template, string) {
		return l.Template.Lookup(name), name
	}, data)
}

// Apply wraps ApplyTemplate with l.Template.
func (l *LazyTemplate) Apply(data interface{}) (*bytes.Buffer, error) {
	return l.ApplyTemplate(nil, data)
}

// ApplyTemplate is internal. Gets a the template, clones and calls BufferExecute.
// getter is called after l.Init so it can rely on l.Template presence.
func (l *LazyTemplate) ApplyTemplate(getter func() (*template.Template, string), data interface{}) (*bytes.Buffer, error) {
	var clone *template.Template
	if err := func() error {
		l.MU.Lock()
		defer l.MU.Unlock()
		if err := l.Init(); err != nil {
			return err
		}
		var t *template.Template
		if getter == nil {
			t = l.Template
		} else {
			var name string
			t, name = getter()
			if t == nil {
				return fmt.Errorf("template: %q is undefined", name)
			}
		}
		var err error
		clone, err = t.Clone()
		return err
	}(); err != nil {
		return nil, err
	}
	return BufferExecute(clone, data)
}

// BufferExecute does t.Execute into buf returned. Does not clone.
func BufferExecute(t *template.Template, data interface{}) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	if err := t.Execute(buf, data); err != nil {
		return nil, err
	}
	return buf, nil
}

// Response is there to .Apply (implicitly) and then .Send explicitly.
// The Send replies with http.Error if it has preceded in Apply.
func (l *LazyTemplate) Response(w http.ResponseWriter, data interface{}) TemplateWriter {
	tw := TemplateWriter{Writer: w}
	tw.Buf, tw.Err = l.Apply(data)
	return tw
}

// TemplateWriter allows applying template first, dealing with error later.
type TemplateWriter struct {
	Writer http.ResponseWriter
	Buf    *bytes.Buffer
	Err    error
}

// Header is a proxy to http.ResponseWriter.Header.
func (tw *TemplateWriter) Header() http.Header {
	if tw.Err != nil {
		return map[string][]string{}
	}
	return tw.Writer.Header()
}

// SetContentLength sets Content-Length header to .Buf length.
func (tw *TemplateWriter) SetContentLength() {
	tw.Header().Set("Content-Length", strconv.Itoa(tw.Buf.Len()))
}

// Send writes to the .Writer what has been buffered in .Buf.
// Or reports the .Err if it has come from LazyTemplate.Response.
func (tw *TemplateWriter) Send() {
	if tw.Err != nil {
		http.Error(tw.Writer, tw.Err.Error(), http.StatusInternalServerError)
	} else {
		io.Copy(tw.Writer, tw.Buf) // or w.Write(Buf.Bytes())
	}
}
