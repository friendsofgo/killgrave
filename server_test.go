package killgrave

import (
	"testing"

	"github.com/gorilla/mux"
)

func TestRunServer(t *testing.T) {
	var serverData = []struct {
		name   string
		server *Server
		err    error
	}{
		{"imposter directory not found", NewServer("failImposterPath", nil), invalidDirectoryError("error")},
		{"malformatted json", NewServer("test/testdata/malformatted", nil), malformattedImposterError("error")},
		{"valid imposter", NewServer("test/testdata/imposters", mux.NewRouter()), nil},
	}

	for _, tt := range serverData {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.server.Run()

			if err == nil {
				if tt.err != nil {
					t.Fatalf("expected an error and got nil")
				}
			}

			if err != nil {
				if tt.err == nil {
					t.Fatalf("not expected any erros and got %+v", err)
				}

				switch err.(type) {
				case invalidDirectoryError:
					if _, ok := (tt.err).(invalidDirectoryError); !ok {
						t.Fatalf("expected invalidDirectoryError got %+v", err)
					}
				case malformattedImposterError:
					if _, ok := (tt.err).(malformattedImposterError); !ok {
						t.Fatalf("expected malformattedImpoasterError got %+v", err)
					}
				default:
					t.Fatalf("not recognize error %+v", err)
				}
			}
		})
	}
}
