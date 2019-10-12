package http

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/gorilla/mux"

	killgrave "github.com/friendsofgo/killgrave/internal"
)

func TestMain(m *testing.M) {
	log.SetOutput(ioutil.Discard)
	os.Exit(m.Run())
}

func TestServer_Build(t *testing.T) {
	var serverData = []struct {
		name   string
		server Server
		err    error
	}{
		{"imposter directory not found", NewServer("failImposterPath", nil, http.Server{}), errors.New("hello")},
		{"malformatted json", NewServer("test/testdata/malformatted_imposters", nil, http.Server{}), nil},
		{"valid imposter", NewServer("test/testdata/imposters", mux.NewRouter(), http.Server{}), nil},
	}

	for _, tt := range serverData {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.server.Build()

			if err == nil {
				if tt.err != nil {
					t.Fatalf("expected an error and got nil")
				}
			}

			if err != nil {
				if tt.err == nil {
					t.Fatalf("not expected any erros and got %+v", err)
				}
			}
		})
	}
}

func TestServer_AccessControl(t *testing.T) {
	config := killgrave.Config{
		ImpostersPath: "imposters",
		Port:          3000,
		Host:          "localhost",
		CORS: killgrave.ConfigCORS{
			Methods:          []string{"GET"},
			Origins:          []string{"*"},
			Headers:          []string{"Content-Type"},
			ExposedHeaders:   []string{"Cache-Control"},
			AllowCredentials: true,
		},
	}

	h := PrepareAccessControl(config.CORS)

	if len(h) <= 0 {
		t.Fatal("Expected any CORS options and got empty")
	}
}
