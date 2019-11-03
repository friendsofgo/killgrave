package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	killgrave "github.com/friendsofgo/killgrave/internal"
	server "github.com/friendsofgo/killgrave/internal/server/http"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/radovskyb/watcher"
)

var (
	_version = "dev"
	_name    = "Killgrave"
)

const (
	_defaultHost          = "localhost"
	_defaultPort          = 3000
	_defaultImpostersPath = "imposters"
	_defaultConfigFile    = ""
)

func main() {
	var (
		host           = flag.String("host", _defaultHost, "if you run your server on a different host")
		port           = flag.Int("port", _defaultPort, "port to run the server")
		imposters      = flag.String("imposters", _defaultImpostersPath, "directory where your imposters are saved")
		showVersion    = flag.Bool("version", false, "show the _version of the application")
		configFilePath = flag.String("config", _defaultConfigFile, "path with configuration file")
		watcherFlag    = flag.Bool("watcher", false, "file watcher, reload the server with each file change")
	)
	flag.Parse()

	if *showVersion {
		fmt.Printf("%s version %s\n", _name, _version)
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

	done := make(chan os.Signal, 1)

	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	srv := runServer(cfg.Host, cfg.Port, cfg)
	srv.Run()

	//Initialize and start the file watcher if the watcher option if true
	w := runWatcher(*watcherFlag, cfg.ImpostersPath, srv, cfg.Host, cfg.Port, cfg)

	<-done
	close(done)
	if w != nil {
		w.Close()
	}

	if err := srv.Shutdown(); err != nil {
		log.Fatal(err)
	}
}

func runWatcher(canWatch bool, pathToWatch string, currentSrv *server.Server, host string, port int, cfg killgrave.Config) *watcher.Watcher {
	if !canWatch {
		return nil
	}
	w, err := killgrave.InitializeWatcher(pathToWatch)
	if err != nil {
		log.Fatal(err)
	}
	killgrave.AttachWatcher(w, func() {
		if err := currentSrv.Shutdown(); err != nil {
			log.Fatal(err)
		}
		currentSrv = runServer(host, port, cfg)
		currentSrv.Run()
	})
	return w
}

func runServer(host string, port int, cfg killgrave.Config) *server.Server {
	router := mux.NewRouter()
	httpAddr := fmt.Sprintf("%s:%d", host, port)

	httpServer := http.Server{
		Addr:    httpAddr,
		Handler: handlers.CORS(server.PrepareAccessControl(cfg.CORS)...)(router),
	}

	s := server.NewServer(
		cfg.ImpostersPath,
		router,
		&httpServer,
	)
	if err := s.Build(); err != nil {
		log.Fatal(err)
	}
	return s
}
