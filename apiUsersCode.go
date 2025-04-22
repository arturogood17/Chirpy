package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id         uuid.UUID `json:"id"`
	Created_at time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_at"`
	Email      string    `json:"email"`
}

func (cfg *apiConfig) hUser(w http.ResponseWriter, req *http.Request) {
	type UserEmail struct {
		Email string `json:"email"`
	}
	decoder := json.NewDecoder(req.Body)
	var userMail UserEmail
	err := decoder.Decode(&userMail)
	if err != nil {
		respondErrorWriter(w, http.StatusInternalServerError, "Error trying to decode de request", err)
		return
	}
	user, err := cfg.Queries.CreateUser(req.Context(), userMail.Email)
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
