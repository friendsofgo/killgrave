package http

import (
	"io"
	"net/http"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

// ImposterHandler create specific handler for the received imposter
func ImposterHandler(i Imposter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res := i.NextResponse()
		if res.Delay.Delay() > 0 {
			time.Sleep(res.Delay.Delay())
		}
		writeHeaders(res, w)
		w.WriteHeader(res.Status)
		writeBody(i, res, w)
		log.WithFields(i.LogFields()).Debugf("Request matched handler")
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

func writeBody(i Imposter, r Response, w http.ResponseWriter) {
	wb := []byte(r.Body)

	if r.BodyFile != nil {
		bodyFile := i.CalculateFilePath(*r.BodyFile)
		wb = fetchBodyFromFile(bodyFile)
	}
	w.Write(wb)
}

func fetchBodyFromFile(bodyFile string) (bytes []byte) {
	if _, err := os.Stat(bodyFile); os.IsNotExist(err) {
		log.Warnf("the body file %s not found\n", bodyFile)
		return
	}

	f, _ := os.Open(bodyFile)
	defer f.Close()
	bytes, err := io.ReadAll(f)
	if err != nil {
		log.Warnf("imposible read the file %s: %v\n", bodyFile, err)
	}
	return
}
