package debugger

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"

	killgrave "github.com/friendsofgo/killgrave/internal"
)

type Debugger struct {
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

func New() *Debugger {
	d := &Debugger{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}

	return d

}

func (d *Debugger) Run() error {
	// TODO: Use a non-default server
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
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

	// TODO: Make debugger address/port configurable
	return http.ListenAndServe(":8080", nil)
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
