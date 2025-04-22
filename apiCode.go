package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id         uuid.UUID `json:"id"`
	Created_at time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_at"`
	Email      string    `json:"email"`
}

func hValChirpy(w http.ResponseWriter, req *http.Request) {
	type chirpyBody struct {
		Body string `json:"body"`
	}
	decoder := json.NewDecoder(req.Body)
	var validC chirpyBody
	if err := decoder.Decode(&validC); err != nil {
		respondErrorWriter(w, http.StatusInternalServerError, "Couldn't decode request body", err)
		return
	}
	new_body := ProfanitiesCleaner(validC.Body)
	const maxlenght = 140
	if len(new_body) > maxlenght {
		respondErrorWriter(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}
	type CleanedB struct {
		CleanedBody string `json:"cleaned_body"`
	}
	responWithJson(w, http.StatusOK, CleanedB{CleanedBody: new_body})
}

func ProfanitiesCleaner(body string) string {
	profanities := map[string]struct{}{
		"kerfuffle": {}, //Los structs vacío se pueden definir así?
		"sharbert":  {},
		"fornax":    {},
	}
	body_split := strings.Fields(body)
	for i, word := range body_split {
		if _, ok := profanities[strings.ToLower(word)]; ok {
			body_split[i] = "****"
		}
	}
	return strings.Join(body_split, " ")
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
