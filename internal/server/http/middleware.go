package http

import (
	"log"
	"net/http"

	"github.com/friendsofgo/killgrave/internal/debugger"
)

type middleware struct {
	handler  http.Handler
	debugger debugger.Debugger
}

func (m *middleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	waitReq, err := m.debugger.NotifyRequestReceived(r)
	if err != nil {
		// TODO: Handle error
		log.Println(err)
	}

	req := waitReq.Wait()
	m.handler.ServeHTTP(w, req)
}
