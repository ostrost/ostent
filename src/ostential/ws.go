package ostential
import (
	"os"
	"log"
	"fmt"
	"time"
	"flag"
	"sync"
	"strings"
	"net/url"
	"net/http"
	"encoding/json"
	gorillawebsocket "github.com/gorilla/websocket"
)

type Duration time.Duration

func(d Duration) String() string {
	s := time.Duration(d).String()
	if strings.HasSuffix(s, "m0s") {
		s = strings.TrimSuffix(s, "0s")
	}
	if strings.HasSuffix(s, "h0m") {
		s = strings.TrimSuffix(s, "0m")
	}
	return s
}

func(d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

type periodValue struct {
	Duration
	above *Duration // optional
}

func(pv *periodValue) Set(input string) error {
	v, err := time.ParseDuration(input)
	if err != nil {
		return err
	}
	if v < time.Second { // hard coded
		return fmt.Errorf("Less than a second: %s", v)
	}
	if v % time.Second != 0 {
		return fmt.Errorf("Not a multiple of a second: %s", v)
	}
	if pv.above != nil && v < time.Duration(*pv.above) {
		return fmt.Errorf("Should be above %s: %s", *pv.above, v)
	}
	pv.Duration = Duration(v)
	return nil
}

var periodFlag = periodValue{Duration: Duration(time.Second)} // default
func init() {
	flag.Var(&periodFlag, "u",      "Collection (update) interval")
	flag.Var(&periodFlag, "update", "Collection (update) interval")
}

var _ = os.Stdout
// var wslog = log.New(os.Stdout, "", log.LstdFlags|log.Lmicroseconds)

func Loop() {
	// flags must be parsed by now, `periodFlag' is used here
	go func() {
		for {
			now := time.Now()
			nextsecond := now.Truncate(time.Second).Add(time.Second).Sub(now)
			<-time.After(nextsecond)

			if connections.expires() {
				// wslog.Printf("Have expires, COLLECT\n")
				lastInfo.collect()
			} else {
				// wslog.Printf("NO REFRESH\n")
			}
			connections.ping()
		}
	}()

	for {
		select {
		case conn := <-register:
			connections.reg(conn)

		case conn := <-unregister:
			close(conn.ping)
			if connections.unreg(conn) == 0 { // if no connections left
				lastInfo.reset_prev()
			}
		}
	}
}

type conn struct {
	*gorillawebsocket.Conn
	ping chan *received
	full client
	mutex sync.Mutex
}

func(c *conn) expires() bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	refreshes := map[string]*refresh{
		// a map by _string_ is for debug only; otherwise that would be a []*refresh
		"expires MEM": c.full.RefreshMEM,
		"expires IF":  c.full.RefreshIF,
		"expires CPU": c.full.RefreshCPU,
		"expires DF":  c.full.RefreshDF,
		"expires PS":  c.full.RefreshPS,
		"expires VG":  c.full.RefreshVG,
	}

	expires := false
	for lprefix, refresh := range refreshes {
		if !refresh.expires() {
			continue
		}
		return true
		log.Println(lprefix, refresh)
		expires = false
	}
	return expires
}

type connmap map[*conn]struct{}
type conns struct {
	connmap
	mutex sync.Mutex
}

func (cs *conns) ping() {
	cs.mutex.Lock()
	defer cs.mutex.Unlock()

	for c := range cs.connmap {
		c.ping <- nil
	}
}

func (cs *conns) reg(c *conn) {
	cs.mutex.Lock()
	defer cs.mutex.Unlock()

	cs.connmap[c] = struct{}{}
}

func (cs *conns) unreg(c *conn) int {
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
	connections = conns{connmap: map[*conn]struct{}{}}
	unregister = make(chan *conn)
	  register = make(chan *conn)
)

type recvClient struct {
	commonClient
	MorePsignal *bool
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

func (rs *recvClient) mergeRefreshSignal(ppinput *string, prefresh *refresh, sendr **refresh, senderr **bool) error {
	if ppinput == nil {
		return nil
	}
	pv := periodValue{above: &periodFlag.Duration}
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

func (rs *recvClient) MergeClient(cs *client, send *sendClient) error {
	rs.mergeMorePsignal(cs)
	if err := rs.mergeRefreshSignal(rs.RefreshSignalMEM, cs.RefreshMEM, &send.RefreshMEM, &send.RefreshErrorMEM); err != nil { return err }
	if err := rs.mergeRefreshSignal(rs.RefreshSignalIF,  cs.RefreshIF,  &send.RefreshIF,  &send.RefreshErrorIF);  err != nil { return err }
	if err := rs.mergeRefreshSignal(rs.RefreshSignalCPU, cs.RefreshCPU, &send.RefreshCPU, &send.RefreshErrorCPU); err != nil { return err }
	if err := rs.mergeRefreshSignal(rs.RefreshSignalDF,  cs.RefreshDF,  &send.RefreshDF,  &send.RefreshErrorDF);  err != nil { return err }
	if err := rs.mergeRefreshSignal(rs.RefreshSignalPS,  cs.RefreshPS,  &send.RefreshPS,  &send.RefreshErrorPS);  err != nil { return err }
	if err := rs.mergeRefreshSignal(rs.RefreshSignalVG,  cs.RefreshVG,  &send.RefreshVG,  &send.RefreshErrorVG);  err != nil { return err }
	return nil
}

type received struct {
	Search *string
	Client *recvClient
}

func(c *conn) writeError(err error) bool {
	type errorJSON struct {
		Error string
	}
	if c.Conn.WriteJSON(errorJSON{err.Error()}) != nil {
		return false
	}
	return true
}

func(c *conn) waitfor_messages() { // read from the conn
	defer c.Conn.Close()
	for {
		rd := new(received)
		if err := c.Conn.ReadJSON(&rd); err != nil {
			break
		}
		c.ping <- rd
	}
}

func(c *conn) waitfor_updates() { // write to the conn
	defer func() {
		unregister <- c
		c.Conn.Close()
	}()
	for {
		select {
		case rd, ok := <- c.ping:
			if !ok {
				break
			}
			if next := c.pong(rd); next != nil {
				if *next {
					continue
				} else {
					break
				}
			}
		}
	}
}

func(c *conn) pong(rd *received) *bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	send := sendClient{}
	var req *http.Request

	if rd != nil {
		if rd.Client != nil {
			err := rd.Client.MergeClient(&c.full, &send)
			if err != nil {
				// if !c.writeError(err) { break }
				// c.writeError(err); continue
				send.DebugError = new(string)
				*send.DebugError = err.Error()
			}
			c.full.Merge(*rd.Client, &send)
		}
		if rd.Search != nil {
			form, err := url.ParseQuery(strings.TrimPrefix(*rd.Search, "?"))
			if err != nil {
				// http.StatusBadRequest?
				// if !c.writeError(err) { break }
				return newtrue()
			}
			req = &http.Request{Form: form}
		}
	}

	updates := getUpdates(req, &c.full, send, false)

	if c.Conn.WriteJSON(updates) != nil {
		return newfalse()
	}
	return nil
}

func slashws(w http.ResponseWriter, req *http.Request) {
	// Upgrader.Upgrade() has Origin check if .CheckOrigin is nil
	upgrader := gorillawebsocket.Upgrader{}
	wsconn, err := upgrader.Upgrade(w, req, nil)
	if err != nil { // Upgrade() does http.Error() to the client
		return
	}

	c := &conn{Conn: wsconn, ping: make(chan *received, 2), full: defaultClient()}
	register <- c
	defer func() {
		unregister <- c
	}()
	go c.waitfor_messages() // read from the client
	   c.waitfor_updates()  // write to the client
}
