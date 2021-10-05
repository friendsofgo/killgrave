package http

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRecorder_Record(t *testing.T) {
	recorder := NewRecorder("test/testdata/recorder/output.imp.json")
	req, err := http.NewRequest(http.MethodGet, "http://localhost/pokemon?limit=100&offset=200", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Trainer", "Ash Ketchum")
	req.Header.Set("Trainer-Key", "25")

	bodyStr := `{"id": 25, name": "Pikachu"}`

	resp := httptest.NewRecorder()
	resp.Body.Write([]byte(bodyStr))
	resp.WriteHeader(http.StatusOK)
	err = recorder.Record(req, resp.Result())

	if err != nil {
		t.Fatal(err)
	}
}
