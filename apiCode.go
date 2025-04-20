package main

import (
	"encoding/json"
	"net/http"
)

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
	const maxlenght = 140
	if len(validC.Body) > maxlenght {
		respondErrorWriter(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}
	type validChirpy struct {
		Valid bool `json:"valid"`
	}
	responWithJson(w, http.StatusOK, validChirpy{Valid: true})
}
