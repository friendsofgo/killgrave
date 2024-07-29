package killgrave

import (
	"errors"
	"io"
	"log"
	"os"
	"testing"
	"time"

	"github.com/radovskyb/watcher"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	log.SetOutput(io.Discard)
	os.Exit(m.Run())
}

func TestAttachWatcher(t *testing.T) {
	var spy bool
	tests := []struct {
		name string
		w    *watcher.Watcher
		fn   func()
	}{
		{"attach watcher and process", watcher.New(), func() { spy = true }},
	}
	for _, tt := range tests {

		AttachWatcher(tt.w, tt.fn)
		tt.w.TriggerEvent(watcher.Create, nil)
		tt.w.Error <- errors.New("some error")
		tt.w.Close()
		time.Sleep(1 * time.Millisecond)
		assert.True(t, spy, "can't read any events")

	}
}

func TestInitializeWatcher(t *testing.T) {

	tests := []struct {
		name        string
		pathToWatch string
		wantWatcher bool
		wantErr     bool
	}{
		{"intialize valid watcher", "test/testdata", true, false},
		{"invalid directory to watch", "<asdddee", false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := InitializeWatcher(tt.pathToWatch)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if tt.wantWatcher {
				assert.NotNil(t, got)
			}
		})
	}
}
