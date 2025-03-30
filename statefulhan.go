package main

import (
	"fmt"
	"net/http"
)

func (a *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a.hits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (a *apiConfig) NumResquests(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte(fmt.Sprintf("Hits: %v", a.hits.Load())))
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(http.StatusText(http.StatusOK)))
}

func (a *apiConfig) Reset(res http.ResponseWriter, req *http.Request) {
	a.hits.Store(0)
	res.WriteHeader(http.StatusOK)
	res.Write([]byte("Hits reset to 0"))
}
