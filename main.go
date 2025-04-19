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
	mux.Handle("/app/", cfg.middlewareConfig(http.StripPrefix("/app", http.FileServer(http.Dir("."))))) //sirve todo lo que están en las carpetas
	mux.HandleFunc("GET /admin/metrics", cfg.handlerServerHits)                                         //se agrego el path admin para uso interno
	mux.HandleFunc("POST /admin/reset", cfg.handlerServerHitsReset)
	mux.HandleFunc("GET /api/healthz", hReadiness)
	err := srvr.ListenAndServe()
	if err != nil {
		log.Fatalf("couldn't listen and serve from server - %v", err)
	}
}
