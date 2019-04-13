package killgrave

import (
	"io/ioutil"
	"net/http"
	"os"
)

func imposterHandler(imposter Imposter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(imposter.Response.Status)
		w.Header().Set("Content-Type", imposter.Response.ContentType)

		f, _ := os.Open(imposter.Response.BodyFile)
		defer f.Close()
		bytes, _ := ioutil.ReadAll(f)

		w.Write(bytes)
	}
}
