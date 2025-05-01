package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

const (
	eventU = "user.upgraded"
)

func (a *apiConfig) WebHook(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}
	var p parameters
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&p); err != nil {
		respondErrorWriter(w, http.StatusInternalServerError, "Error decoding parameters", err)
		return
	}
	if p.Event != eventU {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	_, err := a.Queries.ChirpyRed(req.Context(), p.Data.UserID)
	if err != nil {
		respondErrorWriter(w, http.StatusNotFound, "Error upgrading user to red", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
