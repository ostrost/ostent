package ostent

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
)

// NewErrorLog creates a logger and returns a func to defer.
func NewErrorLog() (*log.Logger, func() error) {
	errlog := logrus.New().Writer()
	return log.New(errlog, "", 0), errlog.Close
}

type Access struct {
	TaggedBin bool
	Logger    *logrus.Logger
	Hosts     struct {
		Mutex sync.Mutex
		Hosts map[string]struct{}
	}
}

func NewAccess(taggedbin bool) *Access {
	logger := logrus.New()
	/* logger.Formatter = &logrus.TextFormatter{
		// DisableTimestamp:true,
		FullTimestamp:   true,
		TimestampFormat: "02/Jan/2006:15:04:05 -0700",
	} // */
	logger.Formatter = &AccessFormatter{}
	a := &Access{
		TaggedBin: taggedbin,
		Logger:    logger,
	}
	a.Hosts.Hosts = make(map[string]struct{})
	return a
}

func (a *Access) PanicHandler(w http.ResponseWriter, req *http.Request, recd interface{}) {
	a.Constructor(PanicHandlerFunc(a.TaggedBin, recd)).ServeHTTP(w, req)
}

func (a *Access) Constructor(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		start := time.Now()
		rsp := &Responding{ResponseWriter: w}
		handler.ServeHTTP(rsp, req)

		if entry := a.Entry(start, *rsp, req); a.TaggedBin && rsp.StatusGood() {
			a.InfoForGood(entry, req)
		} else {
			entry.Info("")
		}
	})
}

func (a *Access) Entry(start time.Time, rsp Responding, req *http.Request) *logrus.Entry {
	since := ZEROTIME.Add(time.Since(start)).Format("5.0000s")
	host := RemoteHost(req)

	uri := req.URL.Path // OR req.RequestURI ??
	if req.Form != nil && len(req.Form) > 0 {
		uri += "?" + req.Form.Encode()
	}
	return a.Logger.WithFields(logrus.Fields{
		"code":      rsp.Status,
		"duration":  since,
		"host":      host,
		"method":    req.Method,
		"proto":     req.Proto,
		"referer":   req.Header.Get("Referer"),
		"size":      rsp.Size,
		"uri":       uri,
		"useragent": req.Header.Get("User-Agent"),
	})
}

func (a *Access) InfoForGood(entry *logrus.Entry, req *http.Request) {
	if host := RemoteHost(req); !a.Seen(host) {
		entry.Data["comment"] = fmt.Sprintf(
			";last info-logged successful request from %s", host)
		entry.Info("")
	} else {
		entry.Debug("")
	}
}

func (a *Access) Seen(host string) bool {
	a.Hosts.Mutex.Lock()
	defer a.Hosts.Mutex.Unlock()
	if _, ok := a.Hosts.Hosts[host]; ok {
		return true
	}
	a.Hosts.Hosts[host] = struct{}{}
	return false
}

func RemoteHost(req *http.Request) string {
	host, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		host = req.RemoteAddr
	}
	return host
}

var ZEROTIME, _ = time.Parse("15:04:05", "00:00:00")

type AccessFormatter struct{}

func (af *AccessFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	f := entry.Data
	comment, okc := f["comment"]
	if !okc {
		comment = ""
	} else if s, ok := comment.(string); ok && s != "" {
		comment = "\t" + s
	}
	echo := func(v interface{}) interface{} {
		if s, ok := v.(string); ok && s == "" {
			return "-"
		}
		return v
	}
	return []byte(fmt.Sprintf("%s - - [%s] %#v %d %d %#v %#v\t;%s%s\n",
		f["host"],
		entry.Time.Format("02/Jan/2006:15:04:05 -0700"),
		fmt.Sprintf("%s %s %s", f["method"], f["uri"], f["proto"]),
		f["code"],
		f["size"],
		echo(f["referer"]),
		echo(f["useragent"]),
		f["duration"],
		comment)), nil
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
