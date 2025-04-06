package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/arturogood17/Chirpy/internal/auth"
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
		respondWithError(res, http.StatusBadRequest, "Empty request", nil)
	}
	type UserCreation struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	var UserData UserCreation
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&UserData); err != nil {
		respondWithError(res, 500, "Error decoding Json", err)
		return
	}
	if UserData.Password == "" {
		respondWithError(res, 400, "User needs a password", nil)
		return
	}
	hashedPass, err := auth.HashPassword(UserData.Password)
	if err != nil {
		respondWithError(res, 500, "Error hashing the password", err)
	}
	user, err := a.dbQueries.CreateUser(context.Background(), database.CreateUserParams{
		Email:          UserData.Email,
		HashedPassword: hashedPass,
	})
	if err != nil {
		respondWithError(res, 500, "Error creating user", err)
		return
	}
	new_user := User{
		ID:        user.ID.String(),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}
	respondWithJson(res, 201, new_user)
}

func (a *apiConfig) HandlerChirps(res http.ResponseWriter, req *http.Request) {
	type reqChirp struct {
		Body string `json:"body"`
	}
	TokenAuth, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(res, 401, "Coukldn't find JWT", err)
		return
	}
	userID, err := auth.ValidateJWT(TokenAuth, a.SECRET)
	if err != nil {
		respondWithError(res, 401, "O este", err)
		return
	}

	decoder := json.NewDecoder(req.Body)
	chirp := reqChirp{}
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
		UserID: userID,
	})

	if err != nil {
		respondWithError(res, 500, "Couldn't create the chirp", err)
		return
	}
	nc := Chirp{
		ID:        new_chirp.ID.String(),
		CreatedAt: new_chirp.CreatedAt,
		UpdatedAt: new_chirp.UpdatedAt,
		Body:      new_chirp.Body,
		UserID:    new_chirp.UserID.String(),
	}
	respondWithJson(res, 201, nc)
}

func (a *apiConfig) AllChirps(res http.ResponseWriter, req *http.Request) {
	chirps, err := a.dbQueries.AllChirps(context.Background())
	if err != nil {
		respondWithError(res, 500, "Error retrieving the all chirps from database", err)
		return
	}
	var slice_chirps []Chirp
	for _, chirp := range chirps {
		nc := Chirp{
			ID:        chirp.ID.String(),
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID.String(),
		}
		slice_chirps = append(slice_chirps, nc)
	}
	respondWithJson(res, 200, slice_chirps)
}

func (a *apiConfig) SingleChirp(res http.ResponseWriter, req *http.Request) {
	pathValue := req.PathValue("chirpID")
	if pathValue == "" {
		respondWithError(res, 404, "Chirp not found", nil)
		return
	}
	chirpID, err := uuid.Parse(pathValue)
	if err != nil {
		respondWithError(res, 500, "Error parsing the chirpID", err)
		return
	}

	chirp, err := a.dbQueries.SingleChirp(context.Background(), chirpID)
	if err != nil {
		respondWithError(res, 500, "Error fetching the desired chirp", err)
		return
	}
	nc := Chirp{
		ID:        chirp.ID.String(),
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID.String(),
	}
	respondWithJson(res, 200, nc)
}

func (a *apiConfig) UserLogin(res http.ResponseWriter, req *http.Request) {
	if req.Body == nil {
		respondWithError(res, http.StatusBadRequest, "Empty request", nil)
	}
	type LoginData struct {
		Password         string `json:"password"`
		Email            string `json:"email"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}

	type response struct {
		User
		Token string `json:"token"`
	}
	var UserData LoginData
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&UserData); err != nil {
		respondWithError(res, 500, "Error decoding Json", err)
		return
	}
	user, err := a.dbQueries.SearchUser(context.Background(), UserData.Email)
	if err != nil {
		respondWithError(res, 401, "Incorrect email or password", err)
		return
	}
	if err = auth.CheckPasswordHash(user.HashedPassword, UserData.Password); err != nil {
		respondWithError(res, 401, "Incorrect email or password", err)
		return
	}
	var token string
	if UserData.ExpiresInSeconds == 0 || UserData.ExpiresInSeconds > 3600 {
		token, err = auth.MakeJWT(user.ID, a.SECRET, time.Hour)
		if err != nil {
			respondWithError(res, 500, "Error creating token", err)
		}
	} else {
		token, err = auth.MakeJWT(user.ID, a.SECRET, time.Duration(UserData.ExpiresInSeconds))
		if err != nil {
			respondWithError(res, 500, "Error creating token", err)
		}
	}

	r := response{
		User: User{
			ID:        user.ID.String(),
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
		Token: token,
	}
	respondWithJson(res, 200, r)
}
