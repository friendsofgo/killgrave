package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"path"

	killgrave "github.com/friendsofgo/killgrave/internal"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var (
	version = "dev"
	name    = "Killgrave"
)

func main() {
	host := flag.String("host", "localhost", "if you run your server on a different host")
	port := flag.Int("port", 3000, "port to run the server")
	imposters := flag.String("imposters", "imposters", "directory where your imposters are saved")
	version := flag.Bool("version", false, "show the version of the application")
	configFilePath := flag.String("config", "", "path with configuration file")

	flag.Parse()

	if *version {
		fmt.Printf("%s version %s\n", name, version)
		return
	}

	var config = killgrave.Config{
		ImpostersPath: *imposters,
		Port:          *port,
		Host:          *host,
	}

	if *configFilePath != "" {
		err := killgrave.ReadConfigFile(*configFilePath, &config)

		if err != nil {
			log.Printf("The config file %s couldn't be read, using default configuration instead\n", *configFilePath)
		} else {
			config.ImpostersPath = path.Join(path.Dir(*configFilePath), config.ImpostersPath)
		}
	}

	r := mux.NewRouter()

	s := killgrave.NewServer(config.ImpostersPath, r)
	if err := s.Build(); err != nil {
		log.Fatal(err)
	}

	httpAddr := fmt.Sprintf("%s:%d", config.Host, config.Port)
	log.Printf("The fake server is on tap now: http://%s:%d\n", config.Host, config.Port)
	log.Fatal(http.ListenAndServe(httpAddr, handlers.CORS(s.AccessControl(config.CORS)...)(r)))
}
