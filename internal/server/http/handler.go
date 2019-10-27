package http

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

// ImposterHandler create specific handler for the received imposter
func ImposterHandler(imposter Imposter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if imposter.Delay() > 0 {
			time.Sleep(imposter.Delay())
		}
		writeHeaders(imposter, w)
		w.WriteHeader(imposter.Response.Status)
		writeBody(imposter, w)
	}
}

func writeHeaders(imposter Imposter, w http.ResponseWriter) {
	if imposter.Response.Headers == nil {
		return
	}

	for key, val := range *imposter.Response.Headers {
		w.Header().Set(key, val)
	}
}

func writeBody(imposter Imposter, w http.ResponseWriter) {
	wb := []byte(imposter.Response.Body)

	if imposter.Response.BodyFile != nil {
		bodyFile := imposter.CalculateFilePath(*imposter.Response.BodyFile)
		wb = fetchBodyFromFile(bodyFile)
	}
	w.Write(wb)
}

func fetchBodyFromFile(bodyFile string) (bytes []byte) {
	if _, err := os.Stat(bodyFile); os.IsNotExist(err) {
		log.Printf("the body file %s not found\n", bodyFile)
		return
	}

	f, _ := os.Open(bodyFile)
	defer f.Close()
	bytes, err := ioutil.ReadAll(f)
	if err != nil {
		log.Printf("imposible read the file %s: %v\n", bodyFile, err)
	}
	return
}
