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
		response := imposter.GetResponse()

		if response.Delay() > 0 {
			time.Sleep(response.Delay())
		}
		writeHeaders(response, w)
		w.WriteHeader(response.Status)
		writeBody(imposter, &response, w)
	}
}

func writeHeaders(response Response, w http.ResponseWriter) {
	if response.Headers == nil {
		return
	}

	for key, val := range *response.Headers {
		w.Header().Set(key, val)
	}
}

func writeBody(imposter Imposter, response *Response, w http.ResponseWriter) {
	wb := []byte(response.Body)

	if response.BodyFile != nil {
		bodyFile := imposter.CalculateFilePath(*response.BodyFile)
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
