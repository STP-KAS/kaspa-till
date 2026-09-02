package main

import (
	"flag"
	"log"
	"net/http"

	"kastill/internal/server"
)

func main() {
	addr := flag.String("addr", ":8082", "listen address")
	flag.Parse()
	s, err := server.New(*addr)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Kaspa Till http://localhost%s — L1 merchant, reserved native stable, no L2, no work-credits", *addr)
	log.Fatal(http.ListenAndServe(*addr, s.Handler()))
}
