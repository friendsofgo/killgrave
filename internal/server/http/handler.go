package http

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"time"

	killgrave "github.com/friendsofgo/killgrave/internal"
)

// ImposterHandler create specific handler for the received imposter
func (s *Server) ImposterHandler(imposter killgrave.Imposter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: Improve or separate debugging logic
		request, err := httputil.DumpRequest(r, true)
		if err != nil {
			// TODO: Handle error
			log.Println(err)
		}
		wait := s.debugger.WaitForRequestContinue(request, imposter)
		wait.Wait()

		bytes, err := json.Marshal(imposter)
		if err != nil {
			// TODO: Handle error
			log.Println(err)
		}
		wait = s.debugger.WaitForImposterContinue(bytes)
		evt := wait.Wait()
		if err := json.Unmarshal(evt.Imposter, &imposter); err != nil {
			// TODO: Handle error
			log.Println(err)
		}

		// TODO: Return the whole response
		wait = s.debugger.WaitForResponseContinue([]byte(imposter.Response.Body))
		evt = wait.Wait()
		// TODO: Be careful because this modifies the response from now on
		imposter.Response.Body = string(evt.Response)

		if imposter.Delay() > 0 {
			time.Sleep(imposter.Delay())
		}
		writeHeaders(imposter, w)
		w.WriteHeader(imposter.Response.Status)
		writeBody(imposter, w)
	}
}

func writeHeaders(imposter killgrave.Imposter, w http.ResponseWriter) {
	if imposter.Response.Headers == nil {
		return
	}

	for key, val := range *imposter.Response.Headers {
		w.Header().Set(key, val)
	}
}

func writeBody(imposter killgrave.Imposter, w http.ResponseWriter) {
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
