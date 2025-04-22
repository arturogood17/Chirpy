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
	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`, a.serverHits.Load())))
}

func (a *apiConfig) handlerServerHitsReset(w http.ResponseWriter, req *http.Request) {
	if a.PLATFORM != "dev" {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Reset is only allowed in dev environment."))
		return
	}
	a.Queries.ResetUsersTable(req.Context())
	a.serverHits.Store(0)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Users table reset"))
}
