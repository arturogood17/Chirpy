package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	hits atomic.Int32
}

func (A *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	A.hits.Add(1)
	return next
}

func main() {
	const port = "8080"
	mux := http.NewServeMux()
	mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir(".")))) //Hay que quitar el app porque queremos servir
	//los archivos que están en el dir actual
	mux.HandleFunc("/healthz", handlerReadiness) //no tienes que crear un directorio para el path
	server := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	fmt.Printf("Serving files from on port: %s\n", port)
	log.Fatal(server.ListenAndServe())
}

func handlerReadiness(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "text/plain; charset=utf-8")
	res.WriteHeader(http.StatusOK)
	res.Write([]byte("OK"))
}
