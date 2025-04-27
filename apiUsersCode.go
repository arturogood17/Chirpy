package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/arturogood17/Chirpy/internal/auth"
	"github.com/arturogood17/Chirpy/internal/database"
	"github.com/google/uuid"
)

type User struct {
	Id         uuid.UUID `json:"id"`
	Created_at time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_at"`
	Email      string    `json:"email"`
	Token      string    `json:"token"`
}

func (cfg *apiConfig) hUser(w http.ResponseWriter, req *http.Request) {
	type Param struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	decoder := json.NewDecoder(req.Body)
	var param Param
	err := decoder.Decode(&param)
	if err != nil {
		respondErrorWriter(w, http.StatusInternalServerError, "Error trying to decode de request", err)
		return
	}
	if param.Password == "" {
		respondErrorWriter(w, http.StatusBadRequest, "Password missing in request", nil)
	}
	hashed_pass, err := auth.HashPassword(param.Password)
	if err != nil {
		respondErrorWriter(w, http.StatusInternalServerError, "error creating hash", err)
	}
	user, err := cfg.Queries.CreateUser(req.Context(), database.CreateUserParams{
		Email:          param.Email,
		HashedPassword: hashed_pass,
	})
	if err != nil {
		respondErrorWriter(w, http.StatusInternalServerError, "Error trying to query the database to create user", err)
		return
	}
	responWithJson(w, http.StatusCreated, User{
		Id:         user.ID,
		Created_at: user.CreatedAt,
		Updated_at: user.UpdatedAt,
		Email:      user.Email,
	})
}
