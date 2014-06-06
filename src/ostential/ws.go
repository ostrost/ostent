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
	above *time.Duration // optional
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
	if v < time.Second { // hard coded
		return fmt.Errorf("Less than a second: %s", v)
	}
	if v % time.Second != 0 {
		return fmt.Errorf("Not a multiple of a second: %s", v)
	}
	if pv.above != nil && v < *pv.above {
		return fmt.Errorf("Should not be bellow %s: %s", *pv.above, v)
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
	fullClient client
}

var (
	 wclients  = make(map[ *wclient ]bool)
	  register = make(chan *wclient)
	unregister = make(chan *wclient)
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

func (rs *recvClient) mergeRefreshSignal(ppinput **string, prefresh *refresh, sendr **refresh, senderr **bool) error {
	if *ppinput == nil {
		return nil
	}
	pv := periodValue{above: &periodFlag.Duration}
	if err := pv.Set(**ppinput); err != nil {
		*senderr = newtrue()
		return err
	}
	*senderr = newfalse()
	prefresh.Duration = pv.Duration
	*sendr = new(refresh)
	(**sendr).Duration = pv.Duration
	*ppinput = nil
	return nil
}

func (rs *recvClient) MergeClient(cs *client, sendc *sendClient) error {
	rs.mergeMorePsignal(cs)
	if err := rs.mergeRefreshSignal(&rs.RefreshSignalMEM, cs.RefreshMEM, &sendc.RefreshMEM, &sendc.RefreshErrorMEM); err != nil { return err }
	if err := rs.mergeRefreshSignal(&rs.RefreshSignalIF,  cs.RefreshIF,  &sendc.RefreshIF,  &sendc.RefreshErrorIF);  err != nil { return err }
	if err := rs.mergeRefreshSignal(&rs.RefreshSignalCPU, cs.RefreshCPU, &sendc.RefreshCPU, &sendc.RefreshErrorCPU); err != nil { return err }
	if err := rs.mergeRefreshSignal(&rs.RefreshSignalDF,  cs.RefreshDF,  &sendc.RefreshDF,  &sendc.RefreshErrorDF);  err != nil { return err }
	if err := rs.mergeRefreshSignal(&rs.RefreshSignalPS,  cs.RefreshPS,  &sendc.RefreshPS,  &sendc.RefreshErrorPS);  err != nil { return err }
	if err := rs.mergeRefreshSignal(&rs.RefreshSignalVG,  cs.RefreshVG,  &sendc.RefreshVG,  &sendc.RefreshErrorVG);  err != nil { return err }
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
			var req *http.Request
			var sendc *sendClient
			if rd != nil {
				if rd.Client != nil {
					sendc = new(sendClient)
					err := rd.Client.MergeClient(&wc.fullClient, sendc)
					if err != nil {
						// if !wc.writeError(err) { break }
						// wc.writeError(err); continue
						sendc.DebugError = new(string)
						*sendc.DebugError = err.Error()
					}
					wc.fullClient.Merge(*rd.Client, sendc)
				}
				if rd.Search != nil {
					form, err := parseSearch(*rd.Search)
					if err != nil {
						// http.StatusBadRequest?
						// if !wc.writeError(err) { break }
						continue
					}
					req = &http.Request{Form: form}
				}
			}

			updates := getUpdates(req, &wc.fullClient, sendc)

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

	wc := &wclient{ws: ws, ping: make(chan *received, 2), fullClient: defaultClient()}
	register <- wc
	defer func() {
		unregister <- wc
	}()
	go wc.waitfor_messages() // read from client
	   wc.waitfor_updates()  // write to  client
}
