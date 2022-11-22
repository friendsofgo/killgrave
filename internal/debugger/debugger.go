package debugger

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"

	killgrave "github.com/friendsofgo/killgrave/internal"
)

type Debugger struct {
	srv *http.Server
	mux *http.ServeMux

	// TODO: Evaluate if more connections are needed
	conn *websocket.Conn
	// TODO: Evaluate Upgrader best practices
	upgrader websocket.Upgrader

	// waits
	// probably one of each per conn
	waitRequestContinue  *debuggerWait
	waitImposterContinue *debuggerWait
	waitResponseContinue *debuggerWait
}

func New(cfg killgrave.ConfigDebugger) (*Debugger, error) {
	mux := http.NewServeMux()

	srv := &http.Server{
		Addr:              cfg.Address,
		ReadTimeout:       1 * time.Second,
		WriteTimeout:      1 * time.Second,
		IdleTimeout:       30 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		Handler:           mux,
	}

	d := &Debugger{
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

//go:embed app/index.html
var index []byte

func (d *Debugger) registerIndex() {
	d.mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write(index)
		w.WriteHeader(http.StatusOK)
	})
}

//go:embed app/static/*
var public embed.FS

func (d *Debugger) registerStatic() error {
	static, err := fs.Sub(public, "app")
	if err != nil {
		return err
	}

	staticFs := http.FileServer(http.FS(static))
	d.mux.Handle("/static/", staticFs)
	return nil
}

func (d *Debugger) registerWs() {
	d.mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		var err error
		d.conn, err = d.upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			// TODO: Return HTTP error response
		}

		// TODO: Evaluate what happens with multiple connections
		// TODO: Manage connection closed/lost
		go func(conn *websocket.Conn) {
			// Continuously read messages from conn
			for {
				msgType, msg, err := conn.ReadMessage()
				if err != nil {
					// TODO: Handle error
					log.Println(err)
					return
				}

				fmt.Println(string(msg))

				switch msgType {
				case websocket.TextMessage:
					// UTF-8 encoded JSON
					d.processWsMessage(msg)

				case websocket.BinaryMessage:
					// TODO: Implement

				case websocket.CloseMessage:
					// TODO: Implement

				default:
					// We just ignore PING/PONG messages
					// and/or those messages with an unknown type.
				}
			}
		}(d.conn)
	})
}

// TODO: Similar for other types?
func (d *Debugger) processWsMessage(msg []byte) {
	var wsMsg Message

	err := json.Unmarshal(msg, &wsMsg)
	if err != nil {
		// TODO: Handle error
		log.Println(err)
		return
	}

	if !wsMsg.Type.IsEventDebugger() {
		log.Println("Unprocessable entity...")
		return
	}

	evt, err := wsMsg.EventDebugger()
	if err != nil {
		// TODO: Handle error
		log.Println(err)
		return
	}

	switch evt.Type {
	case EventDebuggerTypeRequestContinued:
		d.waitRequestContinue.ch <- evt
		close(d.waitRequestContinue.ch)
	case EventDebuggerTypeImposterContinued:
		d.waitImposterContinue.ch <- evt
		close(d.waitImposterContinue.ch)
	case EventDebuggerTypeResponseContinued:
		d.waitResponseContinue.ch <- evt
		close(d.waitResponseContinue.ch)
	default:
		log.Println("Unprocessable event debugger type...")
	}
}

/// More detailed functions

func (d *Debugger) WaitForRequestContinue(request []byte, imposter killgrave.Imposter) Wait {
	d.waitRequestContinue = &debuggerWait{
		ch: make(chan EventDebugger, 1),
	}

	raw, err := json.Marshal(imposter)
	if err != nil {
		// TODO: Handle the error through the debugger
		log.Println(err)
	}

	if err := d.sendDebuggerEvent(EventDebugger{
		Type:     EventDebuggerTypeRequestReceived,
		Request:  request,
		Imposter: raw,
	}); err != nil {
		// TODO: Handle error
		log.Println(err)
	}

	return d.waitRequestContinue
}

func (d *Debugger) WaitForImposterContinue(imposter []byte) Wait {
	d.waitImposterContinue = &debuggerWait{
		ch: make(chan EventDebugger, 1),
	}

	if err := d.sendDebuggerEvent(EventDebugger{
		Type:     EventDebuggerTypeImposterReceived,
		Imposter: imposter,
	}); err != nil {
		// TODO: Handle error
		log.Println(err)
	}

	return d.waitImposterContinue
}

func (d *Debugger) WaitForResponseContinue(response []byte) Wait {
	d.waitResponseContinue = &debuggerWait{
		ch: make(chan EventDebugger, 1),
	}

	if err := d.sendDebuggerEvent(EventDebugger{
		Type:     EventDebuggerTypeResponseReceived,
		Response: response,
	}); err != nil {
		// TODO: Handle error
		log.Println(err)
	}

	return d.waitResponseContinue
}

func (d *Debugger) sendDebuggerEvent(evt EventDebugger) error {
	// TODO: Better management
	if d.conn == nil {
		return nil
	}

	bytes, err := json.Marshal(evt)
	if err != nil {
		return err
	}

	msg := Message{
		Type:    MessageTypeEventDebugger,
		Payload: bytes,
	}

	if err := d.conn.WriteJSON(msg); err != nil {
		return err
	}

	return nil
}
