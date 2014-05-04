package ostential
import (
	"time"
	"flag"
	"net/url"
	"net/http"
	"github.com/gorilla/websocket"
)

var period = time.Second // default
func init() {
	flag.DurationVar(&period, "u",      time.Second, "Collection (update) interval")
	flag.DurationVar(&period, "update", time.Second, "Collection (update) interval")
}

func Loop() {
	// flags must be parsed by now, `period' is used here
	for {
		select {
		case wc := <-register:
			if len(wclients) == 0 {
				collect() // for at least one new client
			}
			wclients[wc] = true

		case wc := <-unregister:
			delete(wclients, wc)
			close(wc.ping)
			if len(wclients) == 0 {
				reset_prev()
			}

		case <-time.After(period):
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
	ws *websocket.Conn
	ping chan *clientState
	form url.Values
	new_search bool
	fullState clientState
}

var (
	 wclients  = make(map[ *wclient ]bool)
	  register = make(chan *wclient)
	unregister = make(chan *wclient)
)

type received struct {
	Search *string
	State *clientState `json:"State"`
}

func(wc *wclient) waitfor_messages() { // read from client
	defer wc.ws.Close()
	for {
		rd := new(received)
		if err := wc.ws.ReadJSON(&rd); err != nil {
			// fmt.Printf("JSON ERR %s\n", err)
			break
		}
		if rd.State != nil {
			wc.fullState.Merge(*rd.State)
		}
		if rd.Search != nil {
			var err error
			wc.form, err = parseSearch(*rd.Search) // (string(data))
			if err != nil {
				// http.StatusBadRequest
				break
			}
			wc.new_search = true
		}
		wc.ping <- rd.State // != nil
	}
}
func(wc *wclient) waitfor_updates() { // write to  client
	defer func() {
		unregister <- wc
		wc.ws.Close()
	}()
	for {
		select {
		case diffState := <- wc.ping:
			// TODO wc.State intros race, need a lock
			/* dState := clientState{}
			if (diffState != nil) {
				dState = *diffState
			} // */
			updates, _, _, _ := getUpdates(&http.Request{Form: wc.form}, wc.new_search, wc.fullState, diffState) // &dState
			wc.new_search = false

			if err := wc.ws.WriteJSON(updates); err != nil {
				break
			}
		}
	}
}

func slashws(w http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if req.Header.Get("Origin") != "http://"+ req.Host {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	ws, err := websocket.Upgrade(w, req, nil, 1024, 1024)
	if err != nil {
		if _, ok := err.(websocket.HandshakeError); ok {
			http.Error(w, "websocket.Upgrade errd", http.StatusBadRequest)
			return
		}
		panic(err)
	}

	wc := &wclient{ws: ws, ping: make(chan *clientState, 1), fullState: defaultClientState()}
	register <- wc
	defer func() {
		unregister <- wc
	}()
	go wc.waitfor_messages() // read from client
	   wc.waitfor_updates()  // write to  client
}
