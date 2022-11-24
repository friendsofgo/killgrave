package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/radovskyb/watcher"

	killgrave "github.com/friendsofgo/killgrave/internal"
	server "github.com/friendsofgo/killgrave/internal/server/http"
)

var (
	_version = "dev"
	_name    = "Killgrave"
)

const (
	_defaultHost          = "localhost"
	_defaultPort          = 3000
	_defaultSecure        = false
	_defaultImpostersPath = "imposters"
	_defaultConfigFile    = ""
	_defaultProxyMode     = killgrave.ProxyNone
	_defaultDebugger      = false
	_defaultDebuggerAddr  = "localhost:3030"
)

func main() {
	var (
		host           = flag.String("host", _defaultHost, "run your server on a different host")
		port           = flag.Int("port", _defaultPort, "run your server on a different port")
		secure         = flag.Bool("secure", _defaultSecure, "run your server using TLS (https)")
		imposters      = flag.String("imposters", _defaultImpostersPath, "directory where imposters are read from")
		showVersion    = flag.Bool("version", false, "show the version of the application")
		configFilePath = flag.String("config", _defaultConfigFile, "path to the configuration file")
		watcherFlag    = flag.Bool("watcher", false, "enable the file watcher, which reloads the server on every file change")
		proxyModeFlag  = flag.String("proxy-mode", _defaultProxyMode.String(), "proxy mode (choose between 'all', 'missing' or 'none')")
		proxyURLFlag   = flag.String("proxy-url", "", "proxy url, use it in combination with proxy-mode")
		debugger       = flag.Bool("debugger", _defaultDebugger, "run your server with the debugger")
		debuggerAddr   = flag.String("debugger-addr", _defaultDebuggerAddr, "debugger address")
	)

	flag.Parse()

	if *showVersion {
		fmt.Printf("%s version %s\n", _name, _version)
		return
	}

	// The config file is mandatory over the flag options
	cfg, err := killgrave.NewConfig(
		*imposters,
		*host,
		*port,
		*secure,
		killgrave.WithProxyConfiguration(*proxyModeFlag, *proxyURLFlag),
		killgrave.WithWatcherConfiguration(*watcherFlag),
		killgrave.WithConfigFile(*configFilePath),
		killgrave.WithDebuggerConfiguration(*debugger, *debuggerAddr),
	)
	if err != nil {
		log.Println(err)
	}

	done := make(chan os.Signal, 1)

	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	var srv server.Server
	srv = runServer(cfg.Host, cfg.Port, cfg)
	srv.Run()

	// Initialize and start the file watcher if the watcher option is true
	w := runWatcher(cfg.Watcher, cfg.ImpostersPath, &srv, cfg.Host, cfg.Port, cfg)

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
		*currentSrv = runServer(host, port, cfg)
		currentSrv.Run()
	})
	return w
}

func runServer(host string, port int, cfg killgrave.Config) server.Server {
	router := mux.NewRouter()
	httpAddr := fmt.Sprintf("%s:%d", host, port)

	httpServer := &http.Server{
		Addr:    httpAddr,
		Handler: handlers.CORS(server.PrepareAccessControl(cfg.CORS)...)(router),
	}

	proxyServer, err := server.NewProxy(cfg.Proxy.URL, cfg.Proxy.Mode)
	if err != nil {
		log.Fatal(err)
	}

	s, err := server.NewServer(
		cfg.ImpostersPath,
		router,
		httpServer,
		proxyServer,
		cfg.Secure,
		cfg.Debugger,
	)
	if err != nil {
		log.Fatal(err)
	}

	if err := s.Build(); err != nil {
		log.Fatal(err)
	}

	return s
}
