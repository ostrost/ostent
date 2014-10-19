package ostent

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	gorillawebsocket "github.com/gorilla/websocket"
	"github.com/ostrost/ostent/types"
)

type backgroundHandler func()

var (
	jobs = struct {
		mutex sync.Mutex
		added []backgroundHandler
	}{}
)

func addBackground(j backgroundHandler) {
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

func init() {
	addBackground(Loop)
}

// Loop is the ostent background job
func Loop() {
	go func() {
		for {
			now := time.Now()
			nextsecond := now.Truncate(time.Second).Add(time.Second).Sub(now)
			<-time.After(nextsecond)

			if Connections.expires() {
				lastInfo.collect()
			}
			Connections.tack()
		}
	}()

	go func() {
		if err := vgwatch(); err != nil { // vagrant
			// ignoring the error
		}
	}()

	for {
		select {
		case update := <-iUPDATES:
			Connections.push(update)

		case conn := <-register:
			Connections.reg(conn)

		case conn := <-unregister:
			if Connections.unreg(conn) == 0 { // if no connections left
				lastInfo.reset_prev()
			}
		}
	}
}

type conn struct {
	Conn *gorillawebsocket.Conn

	requestOrigin *http.Request

	receive    chan *received
	pushch     chan *indexUpdate
	full       client
	minrefresh types.Duration
	access     *logger

	mutex      sync.Mutex
	writemutex sync.Mutex
}

func (c *conn) expires() bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	refreshes := []*refresh{
		c.full.RefreshMEM,
		c.full.RefreshIF,
		c.full.RefreshCPU,
		c.full.RefreshDF,
		c.full.RefreshPS,
		c.full.RefreshVG,
	}

	for _, refresh := range refreshes {
		if refresh.expires() {
			return true
		}
	}
	return false
}

func (c *conn) tack() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.receive <- nil
}

func (c *conn) push(update *indexUpdate) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.pushch <- update
}

type receiver interface {
	tack()
	push(*indexUpdate)
	reload()
	expires() bool
}

type connmap map[receiver]struct{}
type conns struct {
	connmap
	mutex sync.Mutex
}

func (cs *conns) tack() {
	cs.mutex.Lock()
	defer cs.mutex.Unlock()

	for c := range cs.connmap {
		c.tack()
	}
}

// Reload sends reload signal to all the connections, returns false if there were no connections
func (cs *conns) Reload() bool {
	cs.mutex.Lock()
	defer cs.mutex.Unlock()

	var reloaded bool
	for c := range cs.connmap {
		c.reload()
		reloaded = true
	}
	return reloaded
}

func (cs *conns) push(update *indexUpdate) {
	cs.mutex.Lock()
	defer cs.mutex.Unlock()

	for c := range cs.connmap {
		c.push(update)
	}
}

func (cs *conns) reg(r receiver) {
	cs.mutex.Lock()
	defer cs.mutex.Unlock()

	cs.connmap[r] = struct{}{}
}

func (c *conn) closeChans() {
	c.mutex.Lock()
	defer func() {
		defer c.mutex.Unlock()
		if e := recover(); e != nil {
			if err, ok := e.(error); ok {
				fmt.Printf("CLOSE? PANIC %s\n", err.Error())
			} else {
				fmt.Printf("CLOSE? PANIC (UNDESCRIPT) %+v\n", e)
			}
			panic(e)
		}
	}()
	close(c.receive)
	close(c.pushch)
}

// Len return the number of active connections
func (cs *conns) Len() int {
	cs.mutex.Lock()
	defer cs.mutex.Unlock()
	return len(cs.connmap)
}

func (cs *conns) unreg(c *conn) int {
	c.closeChans()

	cs.mutex.Lock()
	defer cs.mutex.Unlock()

	delete(cs.connmap, c)
	return len(cs.connmap)
}

func (cs *conns) expires() bool {
	cs.mutex.Lock()
	defer cs.mutex.Unlock()

	for c := range cs.connmap {
		if c.expires() {
			return true
		}
	}
	return false
}

var (
	// Connections is an instance of unexported conns type to hold
	// active websocket connections. The only method is Reload.
	Connections = conns{connmap: make(map[receiver]struct{})}

	iUPDATES   = make(chan *indexUpdate) // the channel for off-the-clock indexUpdate[s] to push
	unregister = make(chan *conn)
	register   = make(chan *conn)
)

type recvClient struct {
	commonClient
	MorePsignal      *bool
	RefreshSignalMEM *string
	RefreshSignalIF  *string
	RefreshSignalCPU *string
	RefreshSignalDF  *string
	RefreshSignalPS  *string
	RefreshSignalVG  *string
}

func (rs *recvClient) mergeMorePsignal(cs *client) {
	if rs.MorePsignal == nil {
		return
	}
	if *rs.MorePsignal {
		if cs.psLimit < 65536 {
			cs.psLimit *= 2
		}
	} else if cs.psLimit >= 2 {
		cs.psLimit /= 2
	}
	rs.MorePsignal = nil
}

func (rs *recvClient) mergeRefreshSignal(above types.Duration, ppinput *string, prefresh *refresh, sendr **refresh, senderr **bool) error {
	if ppinput == nil {
		return nil
	}
	pv := types.PeriodValue{Above: &above}
	if err := pv.Set(*ppinput); err != nil {
		*senderr = newtrue()
		return err
	}
	*senderr = newfalse()
	*sendr = new(refresh)
	(**sendr).Duration = pv.Duration
	prefresh.Duration = pv.Duration
	prefresh.tick = 0
	return nil
}

func (rs *recvClient) MergeClient(c *conn, send *sendClient) error {
	cs := &c.full
	rs.mergeMorePsignal(cs)
	if err := rs.mergeRefreshSignal(c.minrefresh, rs.RefreshSignalMEM, cs.RefreshMEM, &send.RefreshMEM, &send.RefreshErrorMEM); err != nil {
		return err
	}
	if err := rs.mergeRefreshSignal(c.minrefresh, rs.RefreshSignalIF, cs.RefreshIF, &send.RefreshIF, &send.RefreshErrorIF); err != nil {
		return err
	}
	if err := rs.mergeRefreshSignal(c.minrefresh, rs.RefreshSignalCPU, cs.RefreshCPU, &send.RefreshCPU, &send.RefreshErrorCPU); err != nil {
		return err
	}
	if err := rs.mergeRefreshSignal(c.minrefresh, rs.RefreshSignalDF, cs.RefreshDF, &send.RefreshDF, &send.RefreshErrorDF); err != nil {
		return err
	}
	if err := rs.mergeRefreshSignal(c.minrefresh, rs.RefreshSignalPS, cs.RefreshPS, &send.RefreshPS, &send.RefreshErrorPS); err != nil {
		return err
	}
	if err := rs.mergeRefreshSignal(c.minrefresh, rs.RefreshSignalVG, cs.RefreshVG, &send.RefreshVG, &send.RefreshErrorVG); err != nil {
		return err
	}
	return nil
}

type received struct {
	Search *string
	Client *recvClient
}

type served struct {
	conn     *conn // passing conn into received.ServeHTTP
	received *received
}

func (c *conn) writeJSON(data interface{}) error {
	c.writemutex.Lock()
	defer c.writemutex.Unlock()
	return c.Conn.WriteJSON(data)
}

func (c *conn) reload() {
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
			if next := c.process(rd); next != nil {
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
			if next := c.writeUpdate(*update); !next {
				return
			}
		case _ = <-stop:
			return
		}
	}
}

func (c *conn) process(rd *received) *bool {
	c.mutex.Lock()
	defer func() {
		c.mutex.Unlock()
		if e := recover(); e != nil {
			stack := ""
			/* The stack to be fmt.Printf-d. Not sure if I should
			sbuf := make([]byte, 4096)
			size := runtime.Stack(sbuf, false)
			stack = string(sbuf[:size])
			*/

			if err, ok := e.(error); ok {
				c.writeError(err) // an alert for the client

				fmt.Printf("PANIC %s\n%s\n", err.Error(), stack)
			} else {
				fmt.Printf("PANIC (UNDESCRIPT) %+v\n%s\n", e, stack)
			}
			panic(e)
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
	send := sendClient{}
	if sd.received != nil {
		if sd.received.Client != nil {
			err := sd.received.Client.MergeClient(sd.conn, &send)
			if err != nil {
				// if !sd.conn.Conn.writeError(err) { stop(); return }
				send.DebugError = new(string)
				*send.DebugError = err.Error()
			}
			sd.conn.full.Merge(*sd.received.Client, &send)
		}
	}

	update := getUpdates(r, &sd.conn.full, send, sd.received != nil && sd.received.Client != nil)
	if update == (indexUpdate{}) { // nothing scheduled for the moment, no update
		return
	}

	if sd.conn.writeJSON(update) != nil {
		stop()
		return
	}
	w.WriteHeader(http.StatusSwitchingProtocols) // last change to WriteHeader. 101 is 200
}

func (c *conn) writeUpdate(update indexUpdate) bool {
	if *c.full.HideVG {
		// TODO other .Vagrant* fields may not be discarded
		update.VagrantMachines = nil
	}
	return c.writeJSON(update) == nil
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
	return len(b), nil
}

func SlashwsFunc(access *logger, minrefresh types.Duration) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		Slashws(access, minrefresh, w, req)
	}
}

func Slashws(access *logger, minrefresh types.Duration, w http.ResponseWriter, req *http.Request) {
	// Upgrader.Upgrade() has Origin check if .CheckOrigin is nil
	upgrader := gorillawebsocket.Upgrader{}
	wsconn, err := upgrader.Upgrade(w, req, nil)
	if err != nil { // Upgrade() does http.Error() to the client
		return
	}

	// req.Method == "GET" asserted by the mux
	req.Form = nil // reset reused later .Form
	c := &conn{
		Conn: wsconn,

		requestOrigin: req,

		receive:    make(chan *received, 2),
		pushch:     make(chan *indexUpdate, 2),
		full:       defaultClient(minrefresh),
		minrefresh: minrefresh,
		access:     access,
	}
	register <- c
	defer func() {
		unregister <- c
		c.Conn.Close()
	}()
	stop := make(chan struct{}, 1)
	go c.receiveLoop(stop) // read from the client
	c.updateLoop(stop)     // write to the client
}
