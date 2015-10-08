package ostent

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
)

// NewErrorLog creates a logger and returns a func to defer.
func NewErrorLog() (*log.Logger, func() error) {
	xlog := logrus.New() // into os.Stderr
	xlog.Formatter = &logrus.TextFormatter{FullTimestamp: true}
	// , TimestampFormat: "02/Jan/2006:15:04:05 -0700",
	wlog := xlog.Writer()
	return log.New(wlog, "", 0), wlog.Close
}

type Access struct {
	TaggedBin bool
	Log       *logrus.Logger
	Hosts     struct {
		Mutex sync.Mutex
		Hosts map[string]struct{}
	}
}

func NewAccess(taggedbin bool) *Access {
	a := &Access{TaggedBin: taggedbin}
	a.Hosts.Hosts = make(map[string]struct{})
	a.Log = logrus.New() // into os.Stderr
	a.Log.Formatter = &AccessFormat{}
	return a
}

func (a *Access) Constructor(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		start := time.Now()
		rsp := &Responding{ResponseWriter: w}
		handler.ServeHTTP(rsp, req)
		a.DoLog(start, req, rsp)
	})
}

// DoLog logs the arguments specifics.
func (a *Access) DoLog(start time.Time, req *http.Request, rsp *Responding) {
	since := ZeroTime.Add(time.Since(start))
	host, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		host = req.RemoteAddr
	}
	uri := req.URL.Path // OR req.RequestURI ??
	if req.Form != nil && len(req.Form) > 0 {
		uri += "?" + req.Form.Encode()
	}
	// msg is unused by formatter
	msg, entry := "msg", a.Log.WithFields(logrus.Fields{
		// every value must be string for formatter
		"host":      host,
		"method":    req.Method,
		"uri":       uri,
		"proto":     req.Proto,
		"code":      fmt.Sprintf("%d", rsp.Status),
		"size":      fmt.Sprintf("%d", rsp.Size),
		"referer":   req.Header.Get("Referer"),
		"useragent": req.Header.Get("User-Agent"),
		"duration":  since.Format("5.0000s"),
	})
	if !a.TaggedBin || !rsp.StatusGood() {
		entry.Info(msg)
		return
	}
	if a.Seen(host) {
		entry.Debug(msg)
		return
	}
	entry.Data["comment"] = fmt.Sprintf(
		";last logged successful request from %s", host)
	entry.Info(msg)
}

// Seen returns whether host has been recorded.
// The record is created if it did not exist.
func (a *Access) Seen(host string) bool {
	a.Hosts.Mutex.Lock()
	defer a.Hosts.Mutex.Unlock()
	if _, ok := a.Hosts.Hosts[host]; ok {
		return true
	}
	a.Hosts.Hosts[host] = struct{}{}
	return false
}

// ZeroTime is zero time for formatting duration from it.
var ZeroTime, _ = time.Parse("15:04:05", "00:00:00")

// AccessFormat is a dummy struct with Format method.
type AccessFormat struct{}

// Format conforms to logrus.Formatter interface.
func (af *AccessFormat) Format(entry *logrus.Entry) ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 256))
	WriteField(buf, entry.Data["host"], " - - [")
	WriteField(buf, entry.Time.Format("02/Jan/2006:15:04:05 -0700"), "] ")
	WriteField(buf, entry.Data["method"])
	WriteField(buf, entry.Data["uri"])
	WriteField(buf, entry.Data["proto"])
	WriteField(buf, entry.Data["code"])
	WriteField(buf, entry.Data["size"])
	WriteField(buf, strconv.Quote(entry.Data["referer"].(string)))
	WriteField(buf, strconv.Quote(entry.Data["useragent"].(string)), "\t;")
	WriteField(buf, entry.Data["duration"], "")
	if comment, ok := entry.Data["comment"]; ok {
		WriteField(buf, "\t", comment.(string))
	}
	WriteField(buf, "\n", "")
	return buf.Bytes(), nil
}

// WriteField is a helper to write val and optional posts[0] or a whitespace.
// The name relate to use with logrus.Entry.Data fields, no other Field relation.
func WriteField(buf *bytes.Buffer, val interface{}, posts ...string) {
	if s := val.(string); s == "" {
		buf.WriteString("-")
	} else if s == `""` {
		buf.WriteString(`"-"`)
	} else {
		buf.WriteString(s)
	}
	if len(posts) == 0 {
		buf.WriteString(" ")
	} else { // just [0] is written
		buf.WriteString(posts[0])
	}
}

type Responding struct {
	http.ResponseWriter
	http.Flusher // ?
	Status       int
	Size         int
}

func (r Responding) StatusGood() bool {
	var good = []int{
		http.StatusSwitchingProtocols, // 101
		http.StatusOK,                 // 200
		http.StatusNotModified,        // 304
	}
	for _, g := range good {
		if r.Status == g {
			return true
		}
	}
	return false
}

func (r *Responding) Flush() {
	if f, ok := r.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

func (r *Responding) WriteHeader(status int) {
	r.ResponseWriter.WriteHeader(status)
	r.Status = status
}

func (r *Responding) Write(text []byte) (int, error) {
	if r.Status == 0 { // generic approach to Write-ing before WriteHeader call
		r.WriteHeader(http.StatusOK)
	}
	size, err := r.ResponseWriter.Write(text)
	if err == nil {
		r.Size += size
	}
	return size, err
}

func (r *Responding) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if hj, ok := r.ResponseWriter.(http.Hijacker); ok {
		return hj.Hijack()
	}
	return nil, nil, fmt.Errorf("Responding's ResponseWriter doesn't support the Hijacker interface")
}
