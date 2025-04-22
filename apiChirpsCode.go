package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"errors"

	"github.com/arturogood17/Chirpy/internal/database"
	"github.com/google/uuid"
)

type Chirp struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    string    `json:"user_id"`
}

func (a *apiConfig) hChirp(w http.ResponseWriter, req *http.Request) {
	type param struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}
	decoder := json.NewDecoder(req.Body)
	var paramChirp param
	if err := decoder.Decode(&paramChirp); err != nil {
		respondErrorWriter(w, http.StatusInternalServerError, "Couldn't decode request body", err)
		return
	}
	new_body, err := ProfanitiesCleaner(paramChirp.Body)
	if err != nil {
		respondErrorWriter(w, http.StatusBadRequest, err.Error(), err)
		return
	}

	chirp, err := a.Queries.CreateChirp(req.Context(), database.CreateChirpParams{Body: new_body, UserID: paramChirp.UserID})
	if err != nil {
		respondErrorWriter(w, http.StatusInternalServerError, "Error creating chirp", err)
		return
	}
	responWithJson(w, http.StatusCreated, Chirp{
		ID:        chirp.ID.String(),
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID.String(),
	})

}

func ProfanitiesCleaner(body string) (string, error) {
	const maxlenght = 140
	if len(body) > maxlenght {
		return "", errors.New("The chirp is too long")
	}
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
	return strings.Join(body_split, " "), nil
}

func (a *apiConfig) hListChirps(w http.ResponseWriter, req *http.Request) {
	chirpList, err := a.Queries.ListChirps(req.Context())
	if err != nil {
		respondErrorWriter(w, http.StatusInternalServerError, "Error getting chirp list", err)
		return
	}
	if len(chirpList) == 0 {
		respondErrorWriter(w, http.StatusNotFound, "no chirps found", nil)
		return
	}
	var JsonedList []Chirp
	for _, chirp := range chirpList {
		JsonedList = append(JsonedList, Chirp{
			ID:        chirp.ID.String(),
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID.String(),
		})
	}
	responWithJson(w, http.StatusOK, JsonedList)
}
