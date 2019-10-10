package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	killgrave "github.com/friendsofgo/killgrave/internal"
)

var (
	version = "dev"
	name    = "Killgrave"
)

func main() {
	host := flag.String("host", "localhost", "if you run your server on a different host")
	port := flag.Int("port", 3000, "port to run the server")
	imposters := flag.String("imposters", "imposters", "directory where your imposters are saved")
	showVersion := flag.Bool("version", false, "show the version of the application")
	configFilePath := flag.String("config", "", "path with configuration file")

	flag.Parse()

	if *showVersion {
		fmt.Printf("%s version %s\n", name, version)
		return
	}

	cfg, err := killgrave.NewConfig(
		*imposters,
		*host,
		*port,
		killgrave.WithConfigFile(*configFilePath),
	)
	if err != nil {
		log.Println(err)
	}

	r := mux.NewRouter()

	s := killgrave.NewServer(cfg.ImpostersPath, r)
	if err := s.Build(); err != nil {
		log.Fatal(err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	srv := startServer(cfg.Host, cfg.Port, s, cfg.CORS, r)
	<-done

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
}

func startServer(host string, port int, ks *killgrave.Server, cors killgrave.ConfigCORS, router *mux.Router) *http.Server {
	httpAddr := fmt.Sprintf("%s:%d", host, port)
	log.Printf("The fake server is on tap now: http://%s:%d\n", host, port)

	srv := &http.Server{
		Addr:    httpAddr,
		Handler: handlers.CORS(ks.AccessControl(cors)...)(router),
	}

	go func() {
		err := srv.ListenAndServe()

		if err != http.ErrServerClosed {
			log.Fatal(err)
		} else {
			log.Println("stopping server...")
		}
	}()

	return srv
}
