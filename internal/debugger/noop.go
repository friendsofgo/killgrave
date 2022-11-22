package debugger

import (
	"net/http"

	killgrave "github.com/friendsofgo/killgrave/internal"
)

// NoOp is a debugger implementation that does nothing.
// Use it either for testing purposes or as no op implementation.
type NoOp struct{}

// NewNoOp initializes a new NoOp debugger.
func NewNoOp() NoOp {
	return NoOp{}
}

// NotifyRequestReceived implements the Debugger interface.
func (n NoOp) NotifyRequestReceived(req *http.Request) (*WaitFor[*http.Request], error) {
	ch := make(chan *http.Request, 1)
	ch <- req
	return &WaitFor[*http.Request]{ch: ch}, nil
}

// NotifyImposterMatched implements the Debugger interface.
func (n NoOp) NotifyImposterMatched(imp killgrave.Imposter) (*WaitFor[killgrave.Imposter], error) {
	ch := make(chan killgrave.Imposter, 1)
	ch <- imp
	return &WaitFor[killgrave.Imposter]{ch: ch}, nil
}

// NotifyResponsePrepared implements the Debugger interface.
func (n NoOp) NotifyResponsePrepared(res []byte) (*WaitFor[[]byte], error) {
	ch := make(chan []byte, 1)
	ch <- res
	return &WaitFor[[]byte]{ch: ch}, nil
}
