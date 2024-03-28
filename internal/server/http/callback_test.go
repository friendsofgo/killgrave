package http

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

var iitermutex = new(sync.Mutex)

func TestCallback(t *testing.T) {
	router := mux.NewRouter()
	httpServer := &http.Server{Handler: router}
	imposterFs := NewImposterFS(afero.NewOsFs())
	server := NewServer("test/testdata/imposters", router, httpServer, &Proxy{}, false, imposterFs)

	testCases := map[string]struct {
		imposter           Imposter
		expectedStatusCode int
		expectedError      bool
	}{
		"valid callback": {
			imposter: Imposter{
				Request: Request{
					Method:   "POST",
					Endpoint: "/gophers",
				},
				Callback: &Callback{
					Request: Request{
						Method:   "GET",
						Endpoint: "http://localhost:8080/test",
					},
					Delay: ResponseDelay{
						delay:  2 * int64(time.Second),
						offset: 0,
					},
				},
				Response: Response{
					Status: http.StatusOK,
					Body:   "hello",
				},
			},
			expectedStatusCode: 200,
		},
		"default callback": {
			imposter: Imposter{
				Request: Request{
					Method:   "POST",
					Endpoint: "/gophers",
				},
				Callback: &Callback{
					Request: Request{
						Method:   "GET",
						Endpoint: "http://localhost:8080/test",
					},
					Delay: ResponseDelay{
						delay:  0,
						offset: 0,
					},
				},
				Response: Response{
					Status: http.StatusOK,
					Body:   "hello",
				},
			},
			expectedStatusCode: 200,
		},
	}

	for key, tt := range testCases {
		wg := new(sync.WaitGroup)
		wg.Add(1)
		t.Run(key, func(t *testing.T) {
			defer func() {
				if tt.expectedError {
					rec := recover()
					assert.NotNil(t, rec)
				}
			}()

			go func() {
				router := http.NewServeMux()

				router.HandleFunc("/test", helloHandler(wg))
				http.ListenAndServe(":8080", router)
			}()

			rec := httptest.NewRecorder()
			handler := http.HandlerFunc(ImposterHandler(tt.imposter))

			err := server.Build()
			assert.Nil(t, err)

			// due the caller service were ahead 1 seconds for the default.
			handler.ServeHTTP(rec, httptest.NewRequest(tt.imposter.Request.Method, tt.imposter.Request.Endpoint, nil))
			assert.Equal(t, tt.expectedStatusCode, rec.Code)
		})
		wg.Wait()
	}

	
}

// if this function were called then we knew this function were called from the callback
func helloHandler(wg *sync.WaitGroup) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("hello"))

		wg.Done()
		// done <- struct{}{}
	}
}
