package ostent

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"

	"github.com/ostrost/ostent/flags"
	"github.com/ostrost/ostent/params"
)

type backgroundHandler func()

var (
	jobs = struct {
		mutex sync.Mutex
		added []backgroundHandler
	}{}
)

func AddBackground(j backgroundHandler) {
	jobs.mutex.Lock()
	defer jobs.mutex.Unlock()
	jobs.added = append(jobs.added, j)
}

func RunBackground(defaultDelay flags.Delay) {
	jobs.mutex.Lock()
	defer jobs.mutex.Unlock()
	for _, j := range jobs.added {
		go j()
	}
}

// SleepTilNextSecond sleeps til precisely next second.
func SleepTilNextSecond() {
	now := time.Now()
	nextsecond := now.Truncate(time.Second).Add(time.Second).Sub(now)
	time.Sleep(nextsecond)
}

// CollectLoop is a ostent background job: collect the metrics.
func CollectLoop() {
	go func() {
		for {
			SleepTilNextSecond()
			lastInfo.collect(Machine{})
		}
	}()
}

// ConnectionsLoop is a ostent background job: serve connections.
func ConnectionsLoop() {
	go func() {
		for {
			SleepTilNextSecond()
			Connections.tick()
			Connections.Tack()
		}
	}()

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
	pushch  chan *IndexUpdate
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
}

type connmap map[receiver]struct{}
type conns struct {
	connmap
	mutex sync.Mutex
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
	close(c.pushch)
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
		case update, ok := <-c.pushch:
			if !ok {
				return
			}
			if err := c.writeJSON(update); err != nil {
				return
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
		// if !c.writeError(err) { return newfalse() } // should I write an error?
		return newtrue() // continue receiving
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
		return newfalse() // write failure, stop receiving
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
	update, updated, err := getUpdates(r, sd.conn.para)
	if err != nil || !updated { // nothing scheduled for the moment, no update
		return
	}

	if sd.conn.writeJSON(update) != nil {
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
		pushch:  make(chan *IndexUpdate, 2),
		para:    params.NewParams(sw.MinDelay, sw.MaxDelay),
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

func newfalse() *bool      { return new(bool) }
func newtrue() *bool       { return newbool(true) }
func newbool(v bool) *bool { b := new(bool); *b = v; return b }
