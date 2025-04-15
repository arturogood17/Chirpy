package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	srvr := http.Server{
		Handler: mux,
		Addr:    ":8080",
	}
	err := srvr.ListenAndServe()
	if err != nil {
		log.Fatalf("couldn't listen and serve from server - %v", err)
	}
}
