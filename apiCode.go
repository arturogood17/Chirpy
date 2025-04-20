package main

import (
	"encoding/json"
	"net/http"
	"strings"
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
