package debugger

import (
	"bufio"
	"bytes"
	"embed"
	"encoding/json"
	"io/fs"
	"log"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/gorilla/websocket"

	killgrave "github.com/friendsofgo/killgrave/internal"
)

// Ws is a debugger implementation based on WebSockets.
type Ws struct {
	srv *http.Server
	mux *http.ServeMux

	// TODO: Evaluate if more connections are needed
	conn *websocket.Conn
	// TODO: Evaluate Upgrader best practices
	upgrader websocket.Upgrader

	// Waits
	waitRequest  *WaitFor[*http.Request]
	waitImposter *WaitFor[killgrave.Imposter]
	waitResponse *WaitFor[[]byte]
}

// NewWs initializes a new Ws debugger.
func NewWs(cfg killgrave.ConfigDebugger) (*Ws, error) {
	mux := http.NewServeMux()

	srv := &http.Server{
		Addr:              cfg.Address,
		ReadTimeout:       1 * time.Second,
		WriteTimeout:      1 * time.Second,
		IdleTimeout:       30 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		Handler:           mux,
	}

	d := &Ws{
		srv: srv,
		mux: mux,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}

	// App's index.html
	d.registerIndex()

	// App's static assets (css, js, etc)
	if err := d.registerStatic(); err != nil {
		return nil, err
	}

	// WebSocket endpoint
	d.registerWs()

	// Start server
	go func() {
		log.Printf("The debugger app has been enabled and is available now on: %s\n", cfg.Address)
		err := d.srv.ListenAndServe()
		if err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	return d, nil
}

func (w *Ws) NotifyRequestReceived(req *http.Request) (*WaitFor[*http.Request], error) {
	r, err := httputil.DumpRequest(req, true)
	if err != nil {
		return nil, err
	}

	reqStr := string(r)

	if err := w.sendDebuggerMessage(WsDebuggerMessage{
		Type:    WsDebuggerMessageTypeRequestReceived,
		Request: &reqStr,
	}); err != nil {
		return nil, err
	}

	w.waitRequest = &WaitFor[*http.Request]{ch: make(chan *http.Request)}

	return w.waitRequest, nil
}

func (w *Ws) NotifyImposterMatched(imp killgrave.Imposter) (*WaitFor[killgrave.Imposter], error) {
	i, err := json.Marshal(&imp)
	if err != nil {
		return nil, err
	}

	impStr := string(i)

	if err := w.sendDebuggerMessage(WsDebuggerMessage{
		Type:     WsDebuggerMessageTypeImposterMatched,
		Imposter: &impStr,
	}); err != nil {
		return nil, err
	}

	w.waitImposter = &WaitFor[killgrave.Imposter]{ch: make(chan killgrave.Imposter)}

	return w.waitImposter, nil
}

func (w *Ws) NotifyResponsePrepared(res []byte) (*WaitFor[[]byte], error) {
	resStr := string(res)

	if err := w.sendDebuggerMessage(WsDebuggerMessage{
		Type:     WsDebuggerMessageTypeResponsePrepared,
		Response: &resStr,
	}); err != nil {
		return nil, err
	}

	w.waitResponse = &WaitFor[[]byte]{ch: make(chan []byte)}

	return w.waitResponse, nil
}

func (w *Ws) sendDebuggerMessage(payload WsDebuggerMessage) error {
	// TODO: Better management
	if w.conn == nil {
		return nil
	}

	if err := w.conn.WriteJSON(WsMessage[WsDebuggerMessage]{
		Type:    WsMessageTypeDebugger,
		Payload: payload,
	}); err != nil {
		return err
	}

	return nil
}

//go:embed app/index.html
var index []byte

func (w *Ws) registerIndex() {
	w.mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write(index)
	})
}

//go:embed app/static/*
var public embed.FS

func (w *Ws) registerStatic() error {
	static, err := fs.Sub(public, "app")
	if err != nil {
		return err
	}

	staticFs := http.FileServer(http.FS(static))
	w.mux.Handle("/static/", staticFs)
	return nil
}

func (w *Ws) registerWs() {
	w.mux.HandleFunc("/ws", func(rw http.ResponseWriter, r *http.Request) {
		var err error
		w.conn, err = w.upgrader.Upgrade(rw, r, nil)
		if err != nil {
			log.Println(err)
			// TODO: Return HTTP error response
		}

		// TODO: Evaluate what happens with multiple connections
		// TODO: Manage connection closed/lost
		go func(conn *websocket.Conn) {
			// Continuously read messages from conn
			for {
				msgType, raw, err := conn.ReadMessage()
				if err != nil {
					// TODO: Handle error
					log.Println(err)
					return
				}

				switch msgType {
				case websocket.TextMessage:
					// UTF-8 encoded JSON
					w.processWsMessage(raw)

				case websocket.BinaryMessage:
					// TODO: Implement

				case websocket.CloseMessage:
					// TODO: Implement

				default:
					// We just ignore PING/PONG messages
					// and/or those messages with an unknown type.
				}
			}
		}(w.conn)
	})
}
func (w *Ws) processWsMessage(raw []byte) {
	var msg WsMessage[WsDebuggerMessage]

	err := json.Unmarshal(raw, &msg)
	if err != nil {
		// TODO: Handle error
		log.Println(err)
		return
	}

	if msg.Type != WsMessageTypeDebugger {
		log.Println("Unprocessable entity...")
		return
	}

	switch msg.Payload.Type {
	case WsDebuggerMessageTypeRequestConfirmed:
		w.notifyRequestConfirmed(msg.Payload)
	case WsDebuggerMessageTypeImposterConfirmed:
		w.notifyImposterConfirmed(msg.Payload)
	case WsDebuggerMessageTypeResponseConfirmed:
		w.notifyResponseConfirmed(msg.Payload)
	default:
		// TODO: Handle
		log.Println("Unprocessable event debugger type...")
	}
}

func (w *Ws) notifyRequestConfirmed(msg WsDebuggerMessage) {
	if w.waitRequest == nil {
		return
	}

	req, err := http.ReadRequest(bufio.NewReader(bytes.NewReader([]byte(*msg.Request))))
	if err != nil {
		// TODO: Handle error
		log.Println(err)
	}

	w.waitRequest.ch <- req
	close(w.waitRequest.ch)
	w.waitRequest = nil
}

func (w *Ws) notifyImposterConfirmed(msg WsDebuggerMessage) {
	if w.waitImposter == nil {
		return
	}

	var imp killgrave.Imposter
	err := json.Unmarshal([]byte(*msg.Imposter), &imp)
	if err != nil {
		// TODO: Handle error
		log.Println(err)
	}

	w.waitImposter.ch <- imp
	close(w.waitImposter.ch)
	w.waitImposter = nil
}

func (w *Ws) notifyResponseConfirmed(msg WsDebuggerMessage) {
	if w.waitResponse == nil {
		return
	}

	w.waitResponse.ch <- []byte(*msg.Response)
	close(w.waitResponse.ch)
	w.waitResponse = nil
}
