package main

import (
	"encoding/json"
	"errors"
	"net/http"
)

func hValChirpy(w http.ResponseWriter, req *http.Request) {
	type chirpyBody struct {
		Body string `json:"body"`
	}
	decoder := json.NewDecoder(req.Body)
	var validC chirpyBody
	if err := decoder.Decode(&validC); err != nil {
		respondErrorWriter(w, http.StatusInternalServerError, err)
		return
	}
	if len(validC.Body) > 140 {
		respondErrorWriter(w, http.StatusBadRequest, errors.New("Chirp is too long"))
		return
	}
	type validChirpy struct {
		Valid bool `json:"valid"`
	}
	GoodChirpy := validChirpy{
		Valid: true,
	}
	ValidData, err := json.Marshal(GoodChirpy)
	if err != nil {
		respondErrorWriter(w, http.StatusInternalServerError, err)
		return
	}
	responWithJson(w, http.StatusOK, ValidData)
}
