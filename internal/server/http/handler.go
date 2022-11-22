package http

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	killgrave "github.com/friendsofgo/killgrave/internal"
	"github.com/friendsofgo/killgrave/internal/debugger"
)

// ImposterHandler creates a specific handler for the given imposter
func ImposterHandler(debugger debugger.Debugger, imposter killgrave.Imposter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		waitReq, err := debugger.NotifyRequestReceived(r)
		if err != nil {
			// TODO: Handle error
			log.Println(err)
		}

		// TODO: Implement
		r = waitReq.Wait()

		waitImp, err := debugger.NotifyImposterMatched(imposter)
		if err != nil {
			// TODO: Handle error
			log.Println(err)
		}

		imp := waitImp.Wait()

		waitRes, err := debugger.NotifyResponsePrepared(prepareResponse(imp))
		if err != nil {
			// TODO: Handle error
			log.Println(err)
		}

		resBody := waitRes.Wait()

		if imp.Delay() > 0 {
			time.Sleep(imp.Delay())
		}
		writeHeaders(imp, w)
		w.WriteHeader(imp.Response.Status)
		w.Write(resBody)
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

func prepareResponse(imposter killgrave.Imposter) []byte {
	wb := []byte(imposter.Response.Body)

	if imposter.Response.BodyFile != nil {
		bodyFile := imposter.CalculateFilePath(*imposter.Response.BodyFile)
		wb = fetchBodyFromFile(bodyFile)
	}

	return wb
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
