package ostent

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/Jeffail/gabs"
	"github.com/gorilla/websocket"

	"github.com/ostrost/ostent/internal/plugins/outputs/ostent" // "ostent" output
	"github.com/ostrost/ostent/params"
)

var (
	// connections keep active ws connections registry.
	connections = conns{connmap: make(map[*conn]struct{})}
	// exporting is a "exporting to" list
	exporting = new(struct {
		rwmu sync.RWMutex
		list []exportingItem
	})
)

// UpdateLoop pushes updates to connections.
func UpdateLoop() {
	for {
		up, ok := <-ostent.Updates.Get()
		if ok {
			connections.update(up)
			lastCopy.set(up)
		}
	}
}

var lastCopy = &updateCopy{}

type updateCopy struct {
	mutex sync.Mutex
	up    *ostent.Update
}

func (uc *updateCopy) get() *ostent.Update {
	uc.mutex.Lock()
	defer uc.mutex.Unlock()
	return uc.up
}

func (uc *updateCopy) set(up *ostent.Update) {
	uc.mutex.Lock()
	defer uc.mutex.Unlock()
	uc.up = up
}

type conn struct {
	logger logger
	Conn   *websocket.Conn

	initialRequest *http.Request
	logRequests    bool

	para *params.Params

	mutex      sync.Mutex
	writemutex sync.Mutex
}

type connmap map[*conn]struct{}
type conns struct {
	connmap
	mutex sync.Mutex
}

func (cs *conns) update(up *ostent.Update) {
	cs.mutex.Lock()
	defer cs.mutex.Unlock()

	for c := range cs.connmap {
		c.Process(nil, up)
	}
}

func (cs *conns) reg(c *conn) {
	cs.mutex.Lock()
	defer cs.mutex.Unlock()
	cs.connmap[c] = struct{}{}
}

func (cs *conns) unreg(c *conn) {
	cs.mutex.Lock()
	defer cs.mutex.Unlock()
	delete(cs.connmap, c)
}

func AddExporter(header, text string) {
	e := exporting
	e.rwmu.Lock()
	defer e.rwmu.Unlock()
	e.list = append(e.list, exportingItem{Header: header, Text: text})
}

type exportingItem struct{ Header, Text string }

func exportingCopy() []exportingItem {
	e := exporting
	e.rwmu.RLock()
	defer e.rwmu.RUnlock()
	dup := make([]exportingItem, len(exporting.list))
	copy(dup, e.list)
	return dup
}

type received struct{ Search *string }

type served struct {
	conn      *conn // in
	WriteFail bool  // out
}

func (c *conn) writeJSON(data interface{}) error {
	c.writemutex.Lock()
	defer c.writemutex.Unlock()
	errch := make(chan error, 1)
	go func(ch chan error) { ch <- c.Conn.WriteJSON(data) }(errch)
	select {
	case err := <-errch:
		return err
	case <-time.After(5 * time.Second):
		return fmt.Errorf("timed out (5s)")
	}
}

func (c *conn) writeError(err error) bool {
	return nil == c.writeJSON(struct {
		Error string
	}{err.Error()})
}

func (c *conn) Process(rd *received, up *ostent.Update) bool {
	c.mutex.Lock()
	defer func() {
		c.mutex.Unlock()
		if e := recover(); e != nil {
			c.logger.Println(e)
			if _, ok := e.(websocket.CloseError); !ok {
				c.writeError(fmt.Errorf("%v", e))
			}
		}
	}()

	form, err := rd.form()
	if err != nil {
		// if !c.writeError(err) { return new(bool) } // should I write an error?
		return true // continue receiving
	}

	decoded := form == nil
	ctx := c.initialRequest.Context()
	ctx = context.WithValue(ctx, crequestDecoded, decoded)
	ctx = context.WithValue(ctx, coutputUpdate, up)
	req := c.initialRequest.WithContext(ctx)

	if !decoded {
		form.Set("search", "true")        // identify this type of requests in logs
		req.URL.RawQuery = form.Encode()  // RawQuery as is does not go into logs though
		req.RequestURI = req.URL.String() // the RequestURI goes into logs
	}

	sd := &served{conn: c}
	serve := sd.ServeHTTP
	if !decoded {
		serve = LogHandler(c.logRequests, sd).ServeHTTP
	}
	serve(nil, req)

	return !sd.WriteFail // false on write failure, stop receiving
}

func (rd *received) form() (url.Values, error) {
	if rd == nil || rd.Search == nil {
		return nil, nil
	}
	return url.ParseQuery(strings.TrimPrefix(*rd.Search, "?"))
	// url.ParseQuery should not return a nil url.Values without an error
}

func (sd *served) ServeHTTP(_ http.ResponseWriter, r *http.Request) {
	data, updated, err := Updates(r, sd.conn.para)
	if err != nil || !updated { // nothing scheduled for the moment, no update
		return
	}
	if err := sd.conn.writeJSON(data); err != nil {
		sd.WriteFail = true
	}
}

// IndexWS serves ws updates.
func (sw ServeWS) IndexWS(w http.ResponseWriter, req *http.Request) {
	// Upgrader.Upgrade() has Origin check if .CheckOrigin is nil
	upgrader := &websocket.Upgrader{HandshakeTimeout: 5 * time.Second}
	wsconn, err := upgrader.Upgrade(w, req, nil)
	if err != nil { // Upgrade() does http.Error() to the client
		return
	}

	c := &conn{
		logger: sw.logger,
		Conn:   wsconn,

		initialRequest: req,
		logRequests:    sw.logRequests,

		para: params.NewParams(),
	}
	connections.reg(c)
	defer func() {
		connections.unreg(c)
		c.Conn.Close()
	}()
	for {
		rd := new(received)
		if err := c.Conn.ReadJSON(&rd); err != nil || !c.Process(rd, nil) {
			return
		}
	}
}

func Fetch(keys *params.FetchKeys) error {
	for i := range keys.Values {
		if err := fetchOne(keys.Values[i], keys.Fragments[i]); err != nil {
			return err
		}
	}
	return nil
}

func address(u *url.URL) (string, string, error) {
	switch u.Scheme {
	case "https":
		u.Scheme = "wss"
	case "http":
		u.Scheme = "ws"
	default:
		return "", "", fmt.Errorf("Unknown scheme for WebSocket connection: %s", u.Scheme)
	}
	host, port, err := net.SplitHostPort(u.Host)
	if err != nil {
		if !strings.HasPrefix(err.Error(), "missing port in address") {
			return "", "", err
		}
		if host == "" {
			host = u.Host
		}
	}
	if port == "" {
		switch u.Scheme {
		case "wss":
			port = "443"
		case "ws":
			port = "80"
		}
	}
	return host, port, nil
}

func fetchOne(k params.FetchKey, keys []string) error {
	host, port, err := address(&k.URL)
	if err != nil {
		return err
	}
	search, err := json.Marshal(struct{ Search string }{k.URL.RawQuery})
	if err != nil {
		return err
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
		if err := fetchOnce(wsconn, keys); err != nil {
			return err
		}
		if k.Times == 0 {
			// 0 is the default value, which encodes 1 time pass
			break
		}
	}
	return nil
}

func fetchOnce(wsconn *websocket.Conn, keys []string) error {
	_, message, err := wsconn.ReadMessage()
	if err != nil {
		return err
	}
	jdata, err := gabs.ParseJSON(message)
	if err != nil {
		return err
	}
	one, many := FetchExtract(jdata, keys)
	_ = jdata.Delete("params") // err is ignored (missing  "params" is the only error)
	if one != nil {
		fmt.Println(one.StringIndent("", "  "))
	} else {
		text, err := json.MarshalIndent(many, "", "  ")
		if err != nil {
			return err
		}
		fmt.Printf("%s\n", text)
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
