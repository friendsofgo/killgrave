package debugger

import (
	"net/http"

	killgrave "github.com/friendsofgo/killgrave/internal"
)

// Debugger defines the expected behavior of a Killgrave debugger.
type Debugger interface {
	NotifyRequestReceived(*http.Request) (*WaitFor[*http.Request], error)
	NotifyImposterMatched(killgrave.Imposter) (*WaitFor[killgrave.Imposter], error)
	NotifyResponsePrepared([]byte) (*WaitFor[[]byte], error)
}

// WaitFor represents an async generic container
// that you can Wait on for the inner value.
type WaitFor[T any] struct {
	ch chan T
}

// Wait waits for the inner value to be
// ready and returns it.
func (w *WaitFor[T]) Wait() T {
	return <-w.ch
}
