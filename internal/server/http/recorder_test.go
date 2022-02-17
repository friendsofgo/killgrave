package http

import (
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRecorder_Record(t *testing.T) {
	tests := []struct {
		name string
		recordPath string
	}{
		{
			name: "JSON record",
			recordPath: "test/testdata/recorder/output.imp.json",
		},
		{
			name: "YAML record",
			recordPath: "test/testdata/recorder/output.imp.yaml",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T){
			recorder := NewRecorder(tt.recordPath)
			req, err := http.NewRequest(http.MethodGet, "http://localhost/items?limit=100&offset=200", nil)
			assert.NoError(t, err)

			req.Header.Set("ItemUser", "Conan")
			req.Header.Set("Item-Key", "25")

			bodyStr := `{"id": 25, name": "Umbrella"}`

			resp := ResponseRecorder{
				Status: http.StatusOK,
				Body: bodyStr,
			}

			err = recorder.Record(req, resp)
			assert.NoError(t, err)

			f, err := os.Stat(tt.recordPath)
			assert.NoError(t, err)

			assert.Greater(t, f.Size(), int64(0), "empty file")

			dir := filepath.Dir(tt.recordPath)
			_ = os.Remove(tt.recordPath)
			_ = os.RemoveAll(dir)
		})
	}
}

func TestRecorder_RecordWithInvalidExtension(t *testing.T) {
	recorder := NewRecorder("test/testdata/recorder/output.imp.dist")

	req, err := http.NewRequest(http.MethodGet, "http://localhost/items?limit=100&offset=200", nil)
	assert.NoError(t, err)

	resp := ResponseRecorder{
		Status: http.StatusOK,
		Body: "",
	}

	err = recorder.Record(req, resp)
	assert.Error(t, err)
	assert.Equal(t, err, errUnrecognizedExtension)
}