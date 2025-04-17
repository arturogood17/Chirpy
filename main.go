package main

import (
	"log"
	"net/http"
)

func main() {
	const port = "8080"

	mux := http.NewServeMux()
	srvr := http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}
	mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir("."))))
	mux.HandleFunc("/healthz", hReadiness) //sirve todo lo que están en las carpetas
	err := srvr.ListenAndServe()
	if err != nil {
		log.Fatalf("couldn't listen and serve from server - %v", err)
	}
}
