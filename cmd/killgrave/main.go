package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	killgrave "github.com/friendsofgo/killgrave/internal"
	"github.com/gorilla/mux"
)

var (
	version = "dev"
	name    = "Killgrave"
)

func main() {
	host := flag.String("host", "localhost", "if you run your server on a different host")
	port := flag.Int("port", 3000, "por to run the server")
	imposters := flag.String("imposters", "imposters", "directory where your imposters are saved")
	v := flag.Bool("version", false, "show the version of the application")
	flag.Parse()

	if *v {
		fmt.Printf("%s version %s\n", name, version)
		return
	}

	r := mux.NewRouter()

	s := killgrave.NewServer(*imposters, r)
	if err := s.Build(); err != nil {
		log.Fatal(err)
	}

	httpAddr := fmt.Sprintf("%s:%d", *host, *port)
	log.Printf("The fake server is on tap now: http://%s:%d\n", *host, *port)
	log.Fatal(http.ListenAndServe(httpAddr, r))
}
