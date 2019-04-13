package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/friendsofgo/killgrave"

	"github.com/gorilla/mux"
)

func main() {
	host := flag.String("host", "localhost", "if you run your server on a different host (default: localhost)")
	port := flag.Int("port", 3000, "por to run the server (default: 3000)")
	imposters := flag.String("imposters", "imposters", "directory where your imposter are saved (default: imposters)")
	flag.Parse()

	r := mux.NewRouter()

	s := killgrave.NewServer(*imposters, r)
	if err := s.Run(); err != nil {
		log.Fatal(err)
	}

	httpAddr := fmt.Sprintf(":%d", *port)
	log.Printf("The fake server is on tap now: http://%s:%d\n", *host, *port)
	log.Fatal(http.ListenAndServe(httpAddr, r))
}
