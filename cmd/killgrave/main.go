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
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/radovskyb/watcher"

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
	watcherFlag := flag.Bool("watcher", false, "file watcher, reload the server with each file change")

	flag.Parse()

	if *showVersion {
		fmt.Printf("%s version %s\n", name, version)
		return
	}

	wt := watcher.New()
	wt.SetMaxEvents(1)
	wt.FilterOps(watcher.Rename, watcher.Move, watcher.Create, watcher.Write)
	defer wt.Close()

	var watcherConfig killgrave.ConfigOpt
	if *watcherFlag {
		watcherConfig = killgrave.WithWatcher(wt)
		go func() {
			if err := wt.Start(time.Millisecond * 100); err != nil {
				log.Fatalln(err)
			}
		}()
	}
	cfg, err := killgrave.NewConfig(
		*imposters,
		*host,
		*port,
		killgrave.WithConfigFile(*configFilePath),
		watcherConfig,
	)

	if err != nil {
		log.Println(err)
	}

	exit := make(chan struct{})
	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	var httpSrv http.Server

	httpSrv = createServer(cfg.Host, cfg.Port, cfg)
	runServer(&httpSrv)

	go func() {
		for {
			select {
			case <-done:
				shutDownServer(&httpSrv)
				exit <- struct{}{}
				return
			case evt := <-wt.Event:
				log.Println("modified file:", evt.Name())
				shutDownServer(&httpSrv)
				httpSrv = createServer(cfg.Host, cfg.Port, cfg)
				runServer(&httpSrv)
				time.Sleep(1 * time.Millisecond)
			case err := <-wt.Error:
				log.Printf("Error checking file change: %+v", err)
			default:
				continue
			}
		}
	}()

	<-exit
	close(done)
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

func createServer(host string, port int, cfg killgrave.Config) http.Server {
	r := mux.NewRouter()
	s := killgrave.NewServer(cfg.ImpostersPath, r)
	if err := s.Build(); err != nil {
		log.Fatal(err)
	}

	httpAddr := fmt.Sprintf("%s:%d", host, port)
	log.Printf("The fake server is on tap now: http://%s:%d\n", host, port)

	srv := http.Server{
		Addr:    httpAddr,
		Handler: handlers.CORS(s.AccessControl(cfg.CORS)...)(r),
	}

	return srv
}
