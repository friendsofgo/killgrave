package http

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

// CallbackMap is a map of all the callbacks
type callbackMap map[*time.Ticker]*Callback

var mutex = &sync.Mutex{}

// Callback represent the structure of real callback
// with additional delay
type Callback struct {
	Ticker  *time.Ticker  `json:"-"`
	Request Request       `json:"request" yaml:"request"`
	Delay   ResponseDelay `json:"delay" yaml:"delay"`
}

func (c Callback) Call() {
	log.Println("Initiated Callback")
	buf := new(bytes.Buffer)

	req, err := http.NewRequest(c.Request.Method, c.Request.Endpoint, buf)
	if err != nil {
		return
	}

	if c.Request.Body != nil {
		json.NewEncoder(buf).Encode(c.Request.Body)
	}

	if c.Request.Headers != nil {
		for key, val := range *c.Request.Headers {
			req.Header.Set(key, val)
		}
	}
	req.Header.Set("User-Agent", "Killgrave")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("there's error in your defined callback service: ", err.Error())
		return
	}
	if resp.StatusCode != http.StatusOK {
		log.Printf("Callback service returned non-OK status code: %d\n", resp.StatusCode)

		// Read and log the response body if it's available.
		responseBody, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			log.Println("Error reading response body:", readErr.Error())
		} else {
			log.Printf("Response Body: %s\n", responseBody)
		}

		resp.Body.Close()
		return
	}

	log.Println("Callback completed successfully - ", req.URL.String())
}

// global variable to store all the callbacks
var callbackMapInstance = make(callbackMap)

// remove a callback from the map
func (c callbackMap) Remove(t *time.Ticker) {
	delete(c, t)
}

// add a callback to the map
func (c callbackMap) Add(t *time.Ticker, callback *Callback) {
	mutex.Lock()
	c[t] = callback
	mutex.Unlock()
}

func (s *Server) callbackCron() {
	for {
		for ticker, callback := range callbackMapInstance {
			select {
			case <-ticker.C:
				callback.Call()
				callbackMapInstance.Remove(ticker)
			default:
				continue
			}
		}
	}
}
