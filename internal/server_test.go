package killgrave

import (
	"testing"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

func TestRunServer(t *testing.T) {
	var serverData = []struct {
		name   string
		server *Server
		err    error
	}{
		{"imposter directory not found", NewServer("failImposterPath", nil), errors.New("hello")},
		{"malformatted json", NewServer("test/testdata/malformatted_imposters", nil), nil},
		{"valid imposter", NewServer("test/testdata/imposters", mux.NewRouter()), nil},
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

func TestAccessControl(t *testing.T) {
	s := NewServer("test/testdata/imposters", mux.NewRouter())
	config := Config{
		ImpostersPath: "imposters",
		Port:          3000,
		Host:          "localhost",
		CORS: ConfigCORS{
			Methods:          []string{"GET"},
			Origins:          []string{"*"},
			Headers:          []string{"Content-Type"},
			ExposedHeaders:   []string{"Cache-Control"},
			AllowCredentials: true,
		},
	}

	h := s.AccessControl(config.CORS)

	if len(h) <= 0 {
		t.Fatal("Expected any CORS options and got empty")
	}
}
