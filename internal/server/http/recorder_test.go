package http

import (
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"testing"
)

func TestRecorder_Record(t *testing.T) {
	outputPath := "test/testdata/recorder/output.imp.json"
	recorder := NewRecorder(outputPath)
	req, err := http.NewRequest(http.MethodGet, "http://localhost/items?limit=100&offset=200", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("ItemUser", "Conan")
	req.Header.Set("Item-Key", "25")

	bodyStr := `{"id": 25, name": "Umbrella"}`

	resp := ResponseRecorder{
		Status: http.StatusOK,
		Body: bodyStr,
	}

	err = recorder.Record(req, resp)

	if err != nil {
		t.Fatal(err)
	}

	f, err := os.Stat(outputPath)
	if os.IsNotExist(err) {
		t.Fatal(err)
	}
	if f.Size() <= 0 {
		t.Fatal(errors.New("empty file"))
	}

	dir := filepath.Dir(outputPath)
	os.Remove(outputPath)
	os.RemoveAll(dir)
}
