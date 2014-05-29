package ostential
import (
	"fmt"
	"time"
	"flag"
	"net/url"
	"net/http"
	gorillawebsocket "github.com/gorilla/websocket"
)

type periodValue struct {
	time.Duration
}

func(pv periodValue) String() string { return pv.Duration.String(); }
func(pv *periodValue) Set(input string) error {
	v, err := time.ParseDuration(input)
	if err != nil {
		return err
	}
	if v <= 0 {
		return fmt.Errorf("Negative interval: %s", v)
	}
	if v <= time.Second {
		return fmt.Errorf("Less than a second: %s", v)
	}
	if v % time.Second != 0 {
		return fmt.Errorf("Not a multiple of a second: %s", v)
	}
	pv.Duration = v
	return nil
}

var periodFlag = periodValue{Duration: time.Second} // default
func init() {
	flag.Var(&periodFlag, "u",      "Collection (update) interval")
	flag.Var(&periodFlag, "update", "Collection (update) interval")
}

func Loop() {
	// flags must be parsed by now, `period' is used here
	for {
		now := time.Now()
		next := now.Truncate(periodFlag.Duration).Add(periodFlag.Duration)
		diff := next.Sub(now)
		select {
		case wc := <-register:
			if len(wclients) == 0 {
				collect() // for at least one new client
			}
			wclients[wc] = true

		case wc := <-unregister:
			close(wc.ping)
			delete(wclients, wc)
			if len(wclients) == 0 {
				reset_prev()
			}

		case <-time.After(diff):
			collect()
			for wc := range wclients {
				wc.ping <- nil // false
			}
		}
	}
}

func parseSearch(search string) (url.Values, error) {
	if search != "" && search[0] == '?' {
		search = search[1:]
	}
	return url.ParseQuery(search)
}

type wclient struct {
	ws *gorillawebsocket.Conn
	ping chan *received
	fullState clientState
}

var (
	 wclients  = make(map[ *wclient ]bool)
	  register = make(chan *wclient)
	unregister = make(chan *wclient)
)

type recvState struct {
	clientState
	MorePsignal *bool
}
func (rs *recvState) mergeMorePsignal(cs *clientState) {
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

type received struct {
	Search *string
	State *recvState
}

func(wc *wclient) waitfor_messages() { // read from client
	defer wc.ws.Close()
	for {
		rd := new(received)
		if err := wc.ws.ReadJSON(&rd); err != nil {
			// fmt.Printf("JSON ERR %s\n", err)
			break
		}
		wc.ping <- rd // != nil
	}
}
func(wc *wclient) waitfor_updates() { // write to  client
	defer func() {
		unregister <- wc
		wc.ws.Close()
	}()
	var form url.Values // one per client
	for {
		select {
		case rd := <- wc.ping:
			new_search := false
			var clientdiff *clientState
			if rd != nil {
				if rd.State != nil {
					rd.State.mergeMorePsignal(&wc.fullState)
					clientdiff = &rd.State.clientState
					wc.fullState.Merge(rd.State.clientState, clientdiff)
				}
				if rd.Search != nil {
					var err error
					form, err = parseSearch(*rd.Search)
					if err != nil {
						// http.StatusBadRequest
						continue
					}
					new_search = true
				}
			}

			updates, _, _, _ := getUpdates(&http.Request{Form: form}, new_search, &wc.fullState, clientdiff)

			if wc.ws.WriteJSON(updates) != nil {
				break
			}
		}
	}
}

func slashws(w http.ResponseWriter, req *http.Request) {
	// Upgrader.Upgrade() has Origin check if .CheckOrigin is nil
	upgrader := gorillawebsocket.Upgrader{}
	ws, err := upgrader.Upgrade(w, req, nil)
	if err != nil { // Upgrade() does http.Error() to the client
		return
	}

	wc := &wclient{ws: ws, ping: make(chan *received, 2), fullState: defaultClientState()}
	register <- wc
	defer func() {
		unregister <- wc
	}()
	go wc.waitfor_messages() // read from client
	   wc.waitfor_updates()  // write to  client
}
