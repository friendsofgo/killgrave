package killgrave

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/textproto"
	"os"
	"strings"
)

// ImposterHandler create specific handler for the received imposter
func ImposterHandler(imposter Imposter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := validateHeaders(imposter, r.Header); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		writeHeaders(imposter, w)
		w.WriteHeader(imposter.Response.Status)
		writeBody(imposter, w)
	}
}

func validateHeaders(imposter Imposter, header http.Header) error {
	if imposter.Request.Headers == nil {
		return nil
	}

	for k, v := range *imposter.Request.Headers {
		mimeTypeKey := textproto.CanonicalMIMEHeaderKey(k)
		if _, ok := header[mimeTypeKey]; !ok {
			return fmt.Errorf("the key %s is not specified on header", k)
		}

		if !compareHeaderValues(v, header[mimeTypeKey]) {
			return fmt.Errorf("the key %s expected: %v got:%v", k, v, header[k])
		}
	}

	return nil
}

func compareHeaderValues(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if strings.ToLower(v) != strings.ToLower(b[i]) {
			return false
		}
	}
	return true
}

func writeHeaders(imposter Imposter, w http.ResponseWriter) {
	if imposter.Response.Headers == nil {
		return
	}

	for key, values := range *imposter.Response.Headers {
		for _, val := range values {
			w.Header().Add(key, val)
		}
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
