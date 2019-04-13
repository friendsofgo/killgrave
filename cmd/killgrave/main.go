package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	host := flag.String("host", "localhost", "if you run your server on a different host (default: localhost)")
	port := flag.Int("port", 3000, "por to run the server (default: 3000)")
	flag.Parse()

	r := mux.NewRouter()

	httpAddr := fmt.Sprintf(":%d", *port)
	log.Printf("The fake server is on tap now: http://%s:%d\n", *host, *port)
	log.Fatal(http.ListenAndServe(httpAddr, r))
}
