package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

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

	httpAddr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	log.Printf("The fake server is on tap now: http://%s:%d\n", cfg.Host, cfg.Port)
	log.Fatal(http.ListenAndServe(httpAddr, handlers.CORS(s.AccessControl(cfg.CORS)...)(r)))
}
