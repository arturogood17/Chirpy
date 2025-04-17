package main

import (
	"fmt"
	"net/http"
)

func (a *apiConfig) middlewareConfig(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		a.serverHits.Add(1)
		next.ServeHTTP(w, req) //Esta parte es importante
	})
}

func (a *apiConfig) handlerServerHits(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Hits: %v", a.serverHits.Load())))
}

func (a *apiConfig) handlerServerHitsReset(w http.ResponseWriter, req *http.Request) {
	a.serverHits.Store(0)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hists reset to 0"))
}
