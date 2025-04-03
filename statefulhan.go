package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/arturogood17/Chirpy/internal/database"
	"github.com/google/uuid"
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

func (a *apiConfig) HandlerChirps(res http.ResponseWriter, req *http.Request) {
	type Chirp struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	type resChirp struct {
		ID        string    `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(req.Body)
	chirp := Chirp{}
	if err := decoder.Decode(&chirp); err != nil {
		respondWithError(res, 500, "Error decoding JSON", err)
		return
	}

	if len(chirp.Body) > 140 || len(chirp.Body) == 0 {
		respondWithError(res, http.StatusBadGateway, "Chirp is too long or is empty", nil)
		return
	}
	bodyVal := WordValidation(chirp.Body)

	new_chirp, err := a.dbQueries.CreateChirps(context.Background(), database.CreateChirpsParams{
		Body:   bodyVal,
		UserID: chirp.UserID,
	})

	if err != nil {
		respondWithError(res, 500, "Couldn't create the chirp", err)
		return
	}
	nc := resChirp{
		ID:        new_chirp.ID.String(),
		CreatedAt: new_chirp.CreatedAt,
		UpdatedAt: new_chirp.UpdatedAt,
		Body:      new_chirp.Body,
		UserID:    new_chirp.UserID,
	}
	respondWithJson(res, 201, nc)
}
