package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"errors"

	"github.com/arturogood17/Chirpy/internal/auth"
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
		Body string `json:"body"`
	}
	authToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondErrorWriter(w, http.StatusUnauthorized, err.Error(), err)
		return
	}

	authUser, err := auth.ValidateJWT(authToken, a.SECRET)
	if err != nil {
		respondErrorWriter(w, http.StatusUnauthorized, err.Error(), err)
		return
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

	chirp, err := a.Queries.CreateChirp(req.Context(), database.CreateChirpParams{Body: new_body, UserID: authUser})
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

func (a *apiConfig) hSingleChirp(w http.ResponseWriter, req *http.Request) {
	StringID := req.PathValue("chirpID")
	chirpID, err := uuid.Parse(StringID)
	if err != nil {
		respondErrorWriter(w, http.StatusBadRequest, "Invalid chirp id", err)
		return
	}
	retrievedChirp, err := a.Queries.SingleChirp(req.Context(), chirpID)
	if err != nil {
		respondErrorWriter(w, http.StatusNotFound, "No chirp found", err)
		return
	}
	responWithJson(w, http.StatusOK, Chirp{
		ID:        retrievedChirp.ID.String(),
		CreatedAt: retrievedChirp.CreatedAt,
		UpdatedAt: retrievedChirp.UpdatedAt,
		Body:      retrievedChirp.Body,
		UserID:    retrievedChirp.UserID.String(),
	})
}

func (a *apiConfig) DeleteChirp(w http.ResponseWriter, req *http.Request) {
	StringID := req.PathValue("chirpID") //Primero revisas que el chirpID esté bien
	chirpID, err := uuid.Parse(StringID)
	if err != nil {
		respondErrorWriter(w, http.StatusBadRequest, "Invalid chirp id", err)
		return
	}
	authToken, err := auth.GetBearerToken(req.Header) //Obtienes el access token que viene en el header "Authorization"
	if err != nil {
		respondErrorWriter(w, http.StatusUnauthorized, "token missing or malformed", err)
		return
	}
	userID, err := auth.ValidateJWT(authToken, a.SECRET) //Validas el JWT a ver si el usuario es válido
	if err != nil {
		respondErrorWriter(w, http.StatusForbidden, "User invalid", err)
		return
	}

	chirp, err := a.Queries.SingleChirp(req.Context(), chirpID) //revisas que exista el chirp
	if err != nil {
		respondErrorWriter(w, http.StatusNotFound, "Couldn't find chirp", err)
		return
	}

	if chirp.UserID != userID { //validas que el usuario que está eliminando el chirp es el dueño del chirp
		respondErrorWriter(w, http.StatusForbidden, "You cannot delete this chirp", err)
		return
	}

	err = a.Queries.DeleteChirp(req.Context(), chirpID) //eliminas el chirp
	if err != nil {
		respondErrorWriter(w, http.StatusInternalServerError, "error deleting chirp", err)
		return
	}
	w.WriteHeader(http.StatusNoContent) //respondes con 204 (no content)
}
