package main

import (
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	serverHits atomic.Int32
}

func main() {
	const port = "8080"
	cfg := apiConfig{}
	mux := http.NewServeMux()
	srvr := http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}
	mux.Handle("/app/", cfg.middlewareConfig(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	mux.HandleFunc("/metrics", cfg.handlerServerHits)
	mux.HandleFunc("/reset", cfg.handlerServerHitsReset)
	mux.HandleFunc("/healthz", hReadiness) //sirve todo lo que están en las carpetas
	err := srvr.ListenAndServe()
	if err != nil {
		log.Fatalf("couldn't listen and serve from server - %v", err)
	}
}
