package killgrave

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"

	"github.com/radovskyb/watcher"
)

func TestMain(m *testing.M) {
	log.SetOutput(ioutil.Discard)
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
		if !spy {
			t.Error("can't read any events")
		}

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
			if (err != nil) != tt.wantErr {
				t.Errorf("InitializeWatcher() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantWatcher && got == nil {
				t.Errorf("InitializeWatcher() got = %v, want a pointer watcher", got)
			}
		})
	}
}
