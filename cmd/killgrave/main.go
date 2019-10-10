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
	//watcher := flag.Bool("watcher", false, "file watcher, reload the server with each file change")

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

	watcherCh := make(chan struct{})
	exit := make(chan struct{})
	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	var httpSrv http.Server

	httpSrv = createServer(cfg.Host, cfg.Port, s, cfg.CORS, r)
	runServer(&httpSrv)

	go func() {
		for {
			select {
			case <-done:
				shutDownServer(&httpSrv)
				exit <- struct{}{}
				return
			case <-watcherCh:
				shutDownServer(&httpSrv)
				httpSrv = createServer(cfg.Host, cfg.Port, s, cfg.CORS, r)
				runServer(&httpSrv)
			default:
				continue
			}
		}
	}()

	<-exit
	close(done)
	close(watcherCh)
	close(exit)
}

func shutDownServer(srv *http.Server) {
	log.Println("stopping server...")
	if err := srv.Shutdown(context.Background()); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
}

func runServer(srv *http.Server) {
	go func() {
		err := srv.ListenAndServe()
		if err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
}

func createServer(host string, port int, ks *killgrave.Server, cors killgrave.ConfigCORS, router *mux.Router) http.Server {
	httpAddr := fmt.Sprintf("%s:%d", host, port)
	log.Printf("The fake server is on tap now: http://%s:%d\n", host, port)

	srv := http.Server{
		Addr:    httpAddr,
		Handler: handlers.CORS(ks.AccessControl(cors)...)(router),
	}

	return srv
}
