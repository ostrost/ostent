package ostent

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/Jeffail/gabs"
	"github.com/gorilla/websocket"

	"github.com/ostrost/ostent/params"
)

type backgroundHandler func()

var (
	jobs = struct {
		mutex sync.Mutex
		added []backgroundHandler
	}{}

	// Exporting has "exporting to" list (after init)
	Exporting ExportingList
)

func AddBackground(j backgroundHandler) {
	jobs.mutex.Lock()
	defer jobs.mutex.Unlock()
	jobs.added = append(jobs.added, j)
}

func RunBackground() {
	jobs.mutex.Lock()
	defer jobs.mutex.Unlock()
	for _, j := range jobs.added {
		go j()
	}
}

// NextSecond returns precisely next second in Time.
func NextSecond() time.Time {
	return time.Now().Truncate(time.Second).Add(time.Second)
}

// NextSecondDelta returns precisely next second in Time and Duration.
func NextSecondDelta() (time.Time, time.Duration) {
	now := time.Now()
	when := now.Truncate(time.Second).Add(time.Second)
	return when, when.Sub(now)
}

// CollectLoop is a ostent background job: collect the metrics.
func CollectLoop() {
	for {
		when, sleep := NextSecondDelta()
		time.Sleep(sleep)
		if len(Exporting) != 0 || Connections.NonZero() {
			lastInfo.collect(when, Connections.NonZeroWantProcs())
		}
		Connections.tick()
		Connections.Tack()
	}
}

// ConnectionsLoop is a ostent background job: serve connections.
func ConnectionsLoop() {
	for {
		select {
		case conn := <-Register:
			Connections.Reg(conn)
		case conn := <-Unregister:
			Connections.Unreg(conn)
		}
	}
}

type conn struct {
	Conn     *websocket.Conn
	ErrorLog *log.Logger

	requestOrigin *http.Request

	receive chan *received
	para    *params.Params
	access  *Access

	mutex      sync.Mutex
	writemutex sync.Mutex
}

func (c *conn) Expired() bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.para.Expired()
}

func (c *conn) Tick() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.para.Tick()
}

func (c *conn) WantProcs() bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.para.NonZeroPsn()
}

func (c *conn) Tack() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.receive <- nil
}

type receiver interface {
	Tick()
	Tack()
	Reload()
	Expired() bool
	CloseChans()
	WantProcs() bool
}

type connmap map[receiver]struct{}
type conns struct {
	connmap
	mutex sync.Mutex
}

func (cs *conns) NonZero() bool {
	cs.mutex.Lock()
	defer cs.mutex.Unlock()
	return len(cs.connmap) != 0
}

func (cs *conns) NonZeroWantProcs() bool {
	cs.mutex.Lock()
	defer cs.mutex.Unlock()

	for c := range cs.connmap {
		if c.WantProcs() {
			return true
		}
	}
	return false
}

func (cs *conns) Tack() {
	cs.mutex.Lock()
	defer cs.mutex.Unlock()

	for c := range cs.connmap {
		if c.Expired() {
			c.Tack()
		}
	}
}

// Reload sends reload signal to all the connections, returns false if there were no connections
func (cs *conns) Reload() bool {
	cs.mutex.Lock()
	defer cs.mutex.Unlock()

	var reloaded bool
	for c := range cs.connmap {
		c.Reload()
		reloaded = true
	}
	return reloaded
}

func (cs *conns) Reg(r receiver) {
	cs.mutex.Lock()
	defer cs.mutex.Unlock()
	cs.connmap[r] = struct{}{}
}

func (c *conn) CloseChans() {
	c.mutex.Lock()
	defer func() {
		defer c.mutex.Unlock()
		if e := recover(); e != nil {
			c.ErrorLog.Printf("close error (recovered panic; from CloseChans) %+v\n", e)
		}
	}()
	close(c.receive)
}

func (cs *conns) Unreg(r receiver) {
	r.CloseChans()

	cs.mutex.Lock()
	defer cs.mutex.Unlock()

	delete(cs.connmap, r)
}

func (cs *conns) tick() {
	cs.mutex.Lock()
	defer cs.mutex.Unlock()
	for c := range cs.connmap {
		c.Tick()
	}
}

// ExportingListing keeps "exporting to" list.
type ExportingListing struct {
	RWMutex sync.RWMutex
	ExportingList
}

func (el *ExportingListing) AddExporter(name string, stringer fmt.Stringer) {
	el.RWMutex.Lock()
	el.ExportingList = append(el.ExportingList,
		struct{ Name, Endpoint string }{name, stringer.String()})
	el.RWMutex.Unlock()
}

// ExportingList is a list to be sorted by .Name:
// - Entries come from CLI flags which may be specified in any order
// - That order not preserved since parsing anyway
type ExportingList []struct{ Name, Endpoint string }

func (el ExportingList) Len() int           { return len(el) }
func (el ExportingList) Less(i, j int) bool { return el[i].Name < el[j].Name }
func (el ExportingList) Swap(i, j int)      { el[i], el[j] = el[j], el[i] }

var (
	// Connections is an instance of unexported conns type to hold
	// active websocket connections. The only method is Reload.
	Connections = conns{connmap: make(map[receiver]struct{})}

	Unregister = make(chan receiver)
	Register   = make(chan receiver, 1)
)

type received struct {
	Search *string
}

type served struct {
	conn     *conn // passing conn into received.ServeHTTP
	received *received
}

func (c *conn) writeJSON(data interface{}) error {
	c.writemutex.Lock()
	defer c.writemutex.Unlock()
	errch := make(chan error, 1)
	go func() { errch <- c.Conn.WriteJSON(data) }()
	select {
	case err := <-errch:
		return err
	case <-time.After(5 * time.Second):
		return fmt.Errorf("timed out (5s)")
	}
}

func (c *conn) Reload() {
	c.writeJSON(struct {
		Reload bool
	}{true})
}

func (c *conn) writeError(err error) bool {
	return nil == c.writeJSON(struct {
		Error string
	}{err.Error()})
}

func (c *conn) receiveLoop(stop chan<- struct{}) { // read from the conn
	for {
		rd := new(received)
		if err := c.Conn.ReadJSON(&rd); err != nil {
			stop <- struct{}{}
			return
		}
		c.receive <- rd
	}
}

func (c *conn) updateLoop(stop <-chan struct{}) { // write to the conn
loop:
	for {
		select {
		case rd, ok := <-c.receive:
			if !ok {
				return
			}
			if next := c.Process(rd); next != nil {
				if *next {
					continue loop
				} else {
					return
				}
			}
		case _ = <-stop:
			return
		}
	}
}

func (c *conn) Process(rd *received) *bool {
	c.mutex.Lock()
	defer func() {
		c.mutex.Unlock()
		if e := recover(); e != nil {
			if _, ok := e.(websocket.CloseError); ok {
				c.ErrorLog.Printf("close error (recovered panic; from Proccess) %+v\n", e)
			} else {
				c.ErrorLog.Printf("ws error (recovered panic; sent to client) %+v\n", e)
				c.writeError(fmt.Errorf("%+v", e))
			}
		}
	}()

	var req *http.Request
	if form, err := rd.form(); err != nil {
		// if !c.writeError(err) { return new(bool) } // should I write an error?
		b := new(bool)
		*b = true // true is to continue receiving
		return b
	} else if form != nil {
		// compile an actual Request
		r := *c.requestOrigin
		r.Form = form
		req = &r // http.Request{Form: form}
	}

	sd := served{conn: c, received: rd}
	serve := sd.ServeHTTP // sd.ServeHTTP survives req being nil

	if req != nil && c.access != nil { // the only case when req.Form is not nil
		// a non-nil req is no-go for access anyway
		serve = c.access.Constructor(sd).ServeHTTP
	}

	w := dummyStatus{}
	serve(w, req)

	if w.status == http.StatusBadRequest {
		return new(bool) // write failure, stop receiving
	}
	return nil
}

func (rd *received) form() (url.Values, error) {
	if rd == nil || rd.Search == nil {
		return nil, nil
	}
	return url.ParseQuery(strings.TrimPrefix(*rd.Search, "?"))
	// url.ParseQuery should not return a nil url.Values without an error
}

func (sd served) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	stop := func() {
		w.WriteHeader(http.StatusBadRequest) // well, not a bad request but a write failure
	}
	data, updated, err := Updates(r, sd.conn.para)
	if err != nil || !updated { // nothing scheduled for the moment, no update
		return
	}

	if sd.conn.writeJSON(data) != nil {
		stop()
		return
	}
	w.WriteHeader(http.StatusSwitchingProtocols) // last change to WriteHeader. 101 is 200
}

type dummyStatus struct { // yet another ResponseWriter
	status int
}

func (w dummyStatus) WriteHeader(s int) {
	w.status = s
	// don't expect any actual WriteHeader. This is dummy after all
}

func (w dummyStatus) Header() http.Header {
	panic("dummyStatus.Header: SHOULD NOT BE USED")
	// return w.ResponseWriter.Header()
	// return make(http.Header) // IF TO RETURN ANYTHING, THAT SHOULD BE ONE http.Header PER dummyStatus
}

func (w dummyStatus) Write(b []byte) (int, error) {
	panic("dummyStatus.Write: SHOULD NOT BE USED")
	// return w.ResponseWriter.Write(b)
	// return len(b), nil
}

// IndexWS serves ws updates.
func (sw ServeWS) IndexWS(w http.ResponseWriter, req *http.Request) {
	// Upgrader.Upgrade() has Origin check if .CheckOrigin is nil
	upgrader := websocket.Upgrader{
		HandshakeTimeout: 5 * time.Second,
	}
	wsconn, err := upgrader.Upgrade(w, req, nil)
	if err != nil { // Upgrade() does http.Error() to the client
		return
	}

	// req.Method == "GET" asserted by the mux
	req.Form = nil // reset reused later .Form
	c := &conn{
		Conn:     wsconn,
		ErrorLog: sw.ErrLog,

		requestOrigin: req,

		receive: make(chan *received, 2),
		para:    params.NewParams(sw.DelayBounds),
		access:  sw.Access,
	}
	Register <- c
	defer func() {
		Unregister <- c
		c.Conn.Close()
	}()
	stop := make(chan struct{}, 1)
	go c.receiveLoop(stop) // read from the client
	c.updateLoop(stop)     // write to the client
}

func Fetch(keys *params.FetchKeys) error {
	if len(keys.Values) == 0 {
		if err := keys.Set(""); err != nil {
			return err
		}
	}
	for i := range keys.Values {
		if err := FetchOne(keys.Values[i], keys.Fragments[i]); err != nil {
			return err
		}
	}
	return nil
}

func FetchOne(k params.FetchKey, keys []string) error {
	switch k.URL.Scheme {
	case "https":
		k.URL.Scheme = "wss"
	case "http":
		k.URL.Scheme = "ws"
	default:
		return fmt.Errorf("Unknown scheme for WebSocket connection: %s", k.URL.Scheme)
	}
	search, err := json.Marshal(struct{ Search string }{k.URL.RawQuery})
	if err != nil {
		return err
	}
	host, port, err := net.SplitHostPort(k.URL.Host)
	if err != nil {
		if !strings.HasPrefix(err.Error(), "missing port in address ") {
			return err
		}
		if host == "" {
			host = k.URL.Host
		}
	}
	if port == "" {
		switch k.URL.Scheme {
		case "wss":
			port = "443"
		case "ws":
			port = "80"
		}
	}
	conn, err := net.Dial("tcp", host+":"+port)
	if err != nil {
		return err
	}
	// conn.SetDeadline(time.Now().Add(time.Second))
	headers := http.Header{}
	headers.Set("User-Agent", "ostent/Go-http-client")
	// headers.Set("Host", host)
	//// headers.Set("Origin", "http://"+k.URL.Host+"/")
	k.URL.Fragment = "" // reset the fragment otherwise ws.NewClient fails
	k.URL.Query().Del("times")
	wsconn, _, err := websocket.NewClient(conn, &k.URL, headers, 10, 10) // 4096, 4096)
	if err != nil {
		return fmt.Errorf("%s: %s", k.URL.String(), err)
	}
	if err = wsconn.WriteMessage(websocket.TextMessage, search); err != nil {
		return err
	}

	// k.Times == -1 means non-stop iterations
	for i := 0; k.Times <= 0 || i < k.Times; i++ {
		_, message, err := wsconn.ReadMessage()
		if err != nil {
			return err
		}
		jdata, err := gabs.ParseJSON(message)
		if err != nil {
			return err
		}
		one, many := FetchExtract(jdata, keys)
		jdata.Delete("params")
		if one != nil {
			fmt.Println(one.StringIndent("", "  "))
		} else {
			text, err := json.MarshalIndent(many, "", "  ")
			if err != nil {
				return err
			}
			fmt.Printf("%s\n", text)
		}
		if k.Times == 0 {
			// 0 is the default value, which encodes 1 time pass
			break
		}
	}
	return nil
}

func FetchExtract(jdata *gabs.Container, keys []string) (*gabs.Container, interface{}) {
	if len(keys) == 0 || (len(keys) == 1 && keys[0] == "") {
		return jdata, nil
	}
	if len(keys) == 1 {
		return jdata.Path(keys[0]), nil
	}
	list := make([]interface{}, len(keys))
	for i, key := range keys {
		one, _ := FetchExtract(jdata, []string{key})
		list[i] = one.Data()
	}
	return nil, list
}
