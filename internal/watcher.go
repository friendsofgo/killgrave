package killgrave

import (
	"fmt"
	"log"
	"time"

	"github.com/radovskyb/watcher"
)

// InitializeWatcher initialize a watcher to check for modification on all files
// in the given path to watch
func InitializeWatcher(pathToWatch string) (*watcher.Watcher, error) {
	w := watcher.New()
	w.SetMaxEvents(1)
	w.FilterOps(watcher.Rename, watcher.Move, watcher.Create, watcher.Write)

	if err := w.AddRecursive(pathToWatch); err != nil {
		return nil, fmt.Errorf("%w: error trying to watch change on %s directory", err, pathToWatch)
	}

	return w, nil
}

// AttachWatcher start the watcher, if any error was produced while the starting process the application would crash
// you need to pass a function, this function is the function that will be executed when the watcher
// receive any event the type of defined on the InitializeWatcher function
func AttachWatcher(w *watcher.Watcher, fn func()) {
	go func() {
		if err := w.Start(time.Millisecond * 100); err != nil {
			log.Fatalln(err)
		}
	}()

	readEventsFromWatcher(w, fn)
}

func readEventsFromWatcher(w *watcher.Watcher, fn func()) {
	go func() {
		for {
			select {
			case evt := <-w.Event:
				log.Println("Modified file:", evt.Name())
				fn()
			case err := <-w.Error:
				log.Printf("Error checking file change: %+v", err)
			case <-w.Closed:
				return
			}
		}
	}()
}
