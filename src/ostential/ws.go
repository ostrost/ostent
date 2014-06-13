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

			if wclients.expires() {
				// wslog.Printf("Have expires, COLLECT\n")
				collect()
			} else {
				// wslog.Printf("NO REFRESH\n")
			}
			wclients.ping()
		}
	}()

	for {
		select {
		case wc := <-register:
			wclients.reg(wc)

		case wc := <-unregister:
			close(wc.ping)
			if wclients.unreg(wc) == 0 { // if no clients left
				reset_prev()
			}
		}
	}
}

type wclient struct {
	ws *gorillawebsocket.Conn
	ping chan *received
	fullClient client
	clientMutex sync.Mutex
}

func(wc *wclient) expires() bool {
	wc.clientMutex.Lock()
	defer wc.clientMutex.Unlock()

	refreshes := map[string]*refresh{
		"expires MEM": wc.fullClient.RefreshMEM,
		"expires IF":  wc.fullClient.RefreshIF,
		"expires CPU": wc.fullClient.RefreshCPU,
		"expires DF":  wc.fullClient.RefreshDF,
		"expires PS":  wc.fullClient.RefreshPS,
		"expires VG":  wc.fullClient.RefreshVG,
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

type clientmap map[*wclient]struct{}
type clientreg struct {
	clientmap
	lock sync.Mutex
}

func (cr *clientreg) ping() {
	cr.lock.Lock()
	defer cr.lock.Unlock()
	for wc := range cr.clientmap {
		wc.ping <- nil
	}
}

func (cr *clientreg) reg(wc* wclient) {
	cr.lock.Lock()
	defer cr.lock.Unlock()
	cr.clientmap[wc] = struct{}{}
}

func (cr *clientreg) unreg(wc* wclient) int {
	cr.lock.Lock()
	defer cr.lock.Unlock()
	delete(cr.clientmap, wc)
	return len(cr.clientmap)
}

func (cr *clientreg) expires() bool {
	cr.lock.Lock()
	defer cr.lock.Unlock()
	for wc := range cr.clientmap {
		if wc.expires() {
			return true
		}
	}
	return false
}

var (
	wclients = clientreg{clientmap: clientmap{}}
	unregister = make(chan *wclient)
	  register = make(chan *wclient)
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

func(wc *wclient) writeError(err error) bool {
	type errorJSON struct {
		Error string
	}
	if wc.ws.WriteJSON(errorJSON{err.Error()}) != nil {
		return false
	}
	return true
}

func(wc *wclient) waitfor_messages() { // read from the client
	defer wc.ws.Close()
	for {
		rd := new(received)
		if err := wc.ws.ReadJSON(&rd); err != nil {
			break
		}
		wc.ping <- rd
	}
}

func(wc *wclient) waitfor_updates() { // write to the client
	defer func() {
		unregister <- wc
		wc.ws.Close()
	}()
	for {
		select {
		case rd, ok := <- wc.ping:
			if !ok {
				break
			}
			if next := wc.pong(rd); next != nil {
				if *next {
					continue
				} else {
					break
				}
			}
		}
	}
}

func(wc *wclient) pong(rd *received) *bool {
	wc.clientMutex.Lock()
	defer wc.clientMutex.Unlock()

	send := sendClient{}
	var req *http.Request

	if rd != nil {
		if rd.Client != nil {
			err := rd.Client.MergeClient(&wc.fullClient, &send)
			if err != nil {
				// if !wc.writeError(err) { break }
				// wc.writeError(err); continue
				send.DebugError = new(string)
				*send.DebugError = err.Error()
			}
			wc.fullClient.Merge(*rd.Client, &send)
		}
		if rd.Search != nil {
			form, err := url.ParseQuery(strings.TrimPrefix(*rd.Search, "?"))
			if err != nil {
				// http.StatusBadRequest?
				// if !wc.writeError(err) { break }
				return newtrue()
			}
			req = &http.Request{Form: form}
		}
	}

	updates := getUpdates(req, &wc.fullClient, send, false)

	if wc.ws.WriteJSON(updates) != nil {
		return newfalse()
	}
	return nil
}

func slashws(w http.ResponseWriter, req *http.Request) {
	// Upgrader.Upgrade() has Origin check if .CheckOrigin is nil
	upgrader := gorillawebsocket.Upgrader{}
	ws, err := upgrader.Upgrade(w, req, nil)
	if err != nil { // Upgrade() does http.Error() to the client
		return
	}

	wc := &wclient{ws: ws, ping: make(chan *received, 2), fullClient: defaultClient()}
	register <- wc
	defer func() {
		unregister <- wc
	}()
	go wc.waitfor_messages() // read from the client
	   wc.waitfor_updates()  // write to the client
}
