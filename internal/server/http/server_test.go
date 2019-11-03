package http

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"

	killgrave "github.com/friendsofgo/killgrave/internal"
	"github.com/gorilla/mux"
)

func TestMain(m *testing.M) {
	log.SetOutput(ioutil.Discard)
	os.Exit(m.Run())
}

func TestServer_Build(t *testing.T) {
	type check func(err error, t *testing.T)

	hasNoError := func() check {
		return func(err error, t *testing.T) {
			t.Helper()
			if err != nil {
				t.Fatalf("not expected any erros and got %+v", err)
			}
		}
	}
	hasNotExistError := func() check {
		return func(err error, t *testing.T) {
			t.Helper()
			if !errors.Is(err, os.ErrNotExist) {
				t.Fatalf("expected to get not exist error, got: %+v", err)
			}
		}
	}
	hasMalformedFileError := func() check {
		return func(err error, t *testing.T) {
			t.Helper()
			if !errors.Is(err, errMalformedFile) {
				t.Fatalf("expected to get malformed file error, got: %+v", err)
			}
		}
	}

	var serverData = []struct {
		name     string
		server   *Server
		errCheck check
	}{
		{
			name:     "imposter directory not found",
			server:   NewServer("failImposterPath", nil, &http.Server{}),
			errCheck: hasNotExistError(),
		},
		{
			name:     "malformatted json",
			server:   NewServer("test/testdata/malformatted_imposters", nil, &http.Server{}),
			errCheck: hasMalformedFileError(),
		},
		{
			name:     "valid imposter",
			server:   NewServer("test/testdata/imposters", mux.NewRouter(), &http.Server{}),
			errCheck: hasNoError(),
		},
	}

	for _, tt := range serverData {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.server.Build()
			tt.errCheck(err, t)
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
