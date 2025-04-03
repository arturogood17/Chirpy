package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
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
	if a.PLATFORM != "dev" {
		log.Fatal("403 Forbidden")
	}
	if err := a.dbQueries.DeleteUsers(context.Background()); err != nil {
		log.Fatalf("error deleting the users from the database")
	}
}

func (a *apiConfig) HandlerUser(res http.ResponseWriter, req *http.Request) {
	if req.Body == nil {
		respondWithError(res, http.StatusBadRequest, "An email is needed to create an user", nil)
	}
	type UserCreation struct {
		Email string `json:"email"`
	}
	var email UserCreation
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&email); err != nil {
		respondWithError(res, 500, "Error decoding Json", err)
		return
	}
	user, err := a.dbQueries.CreateUser(context.Background(), email.Email)
	if err != nil {
		respondWithError(res, 500, "Error creating user", err)
		return
	}
	new_user := User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}
	respondWithJson(res, 201, new_user)
}

func HandlerChirps(res http.ResponseWriter, req *http.Request) {
	type JsonBody struct {
		Body string `json:"body"`
	}
	type ValidR struct {
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(req.Body)
	jb := JsonBody{}
	if err := decoder.Decode(&jb); err != nil {
		respondWithError(res, 500, "Error decoding JSON", err)
		return
	}

	if len(jb.Body) > 140 || len(jb.Body) == 0 {
		respondWithError(res, http.StatusBadGateway, "Chirp is too long or is empty", nil)
		return
	}
	bodyVal := WordValidation(jb.Body)
	valid := ValidR{
		CleanedBody: bodyVal,
	}
	respondWithJson(res, 200, valid)
}
