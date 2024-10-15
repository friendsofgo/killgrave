package http

import (
	"log"
	"net/http"
	"time"
)

// Handler create specific handler for the received imposter
func ImposterHandler(i Imposter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res := i.NextResponse()
		if res.Delay.Delay() > 0 {
			time.Sleep(res.Delay.Delay())
		}
		writeHeaders(res, w)
		w.WriteHeader(res.Status)
		writeBody(res, w)
	}
}

func writeHeaders(r Response, w http.ResponseWriter) {
	if r.Headers == nil {
		return
	}

	for key, val := range *r.Headers {
		w.Header().Set(key, val)
	}
}

func writeBody(r Response, w http.ResponseWriter) {
	_, err := w.Write(r.BodyData)
	if err != nil {
		log.Printf("error writing body: %v\n", err)
	}
}
