package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (a *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a.hits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (a *apiConfig) HandlerMetrics(res http.ResponseWriter, req *http.Request) {
	req.Header.Set("Content-Type", "text/html")
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(fmt.Sprintf(`<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>
`, a.hits.Load())))
}

func (a *apiConfig) HandlerReset(res http.ResponseWriter, req *http.Request) {
	a.hits.Store(0)
	res.WriteHeader(http.StatusOK)
	res.Write([]byte("Hits reset to 0"))
}

func (a *apiConfig) HandlerChirps(res http.ResponseWriter, req *http.Request) {
	type JsonBody struct {
		Body string `json:"body"`
	}
	type ErrorChirp struct {
		Error string `json:"error"`
	}
	type ValidR struct {
		Valid bool `json:"valid"`
	}
	decoder := json.NewDecoder(req.Body)
	jb := JsonBody{}
	if err := decoder.Decode(&jb); err != nil {
		res.WriteHeader(500)
		fmt.Printf("error al decodificar el request body: %v", err)
		return
	}
	if len(jb.Body) > 140 || len(jb.Body) == 0 {
		er := ErrorChirp{
			Error: "Chirp is too long or is empty",
		}
		enc, err := json.Marshal(er)
		if err != nil {
			res.WriteHeader(500)
			fmt.Printf("Something went wrong: %v", err)
			return
		}
		res.WriteHeader(400)
		res.Write(enc)
		return
	}
	valid := ValidR{
		Valid: true,
	}
	val, err := json.Marshal(valid)
	if err != nil {
		res.WriteHeader(500)
		return
	}
	res.WriteHeader(200)
	res.Write(val)
}
