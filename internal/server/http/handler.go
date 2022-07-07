package http

import (
	"github.com/open-telemetry/opamp-go/protobufs"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

// ImposterHandler create specific handler for the received imposter
func ImposterHandler(imposter Imposter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("request headers: %v", r.Header)

		if imposter.Delay() > 0 {
			time.Sleep(imposter.Delay())
		}

		writeHeaders(imposter, w)
		w.WriteHeader(imposter.Response.Status)

		if w.Header().Get("Content-Type") == "application/x-protobuf" {
			writeOpAmpServerToAgentProtoBodyHack(imposter, w)
		} else {
			writeBody(imposter, w)
		}
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

func writeOpAmpServerToAgentProtoBodyHack(imposter Imposter, w http.ResponseWriter) {
	log.Printf("Wring ServerToAgent Protobuf from Payload")
	wb := []byte(imposter.Response.Body)
	if imposter.Response.BodyFile != nil {
		bodyFile := imposter.CalculateFilePath(*imposter.Response.BodyFile)
		wb = fetchBodyFromFile(bodyFile)
	}
	var msg protobufs.ServerToAgent
	err := protojson.Unmarshal(wb, &msg)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	data, _ := proto.Marshal(&msg)
	w.Write(data)

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
