package main

import (
	"encoding/json"
	"net/http"

	"github.com/arturogood17/Chirpy/internal/auth"
	"github.com/arturogood17/Chirpy/internal/database"
)

func (a *apiConfig) Authentication(w http.ResponseWriter, req *http.Request) {
	type param struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	authToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondErrorWriter(w, http.StatusUnauthorized, "Token malformed or missing", err)
		return
	}
	user, err := auth.ValidateJWT(authToken, a.SECRET)
	if err != nil {
		respondErrorWriter(w, http.StatusUnauthorized, "Unathorized", err)
		return
	}

	decoder := json.NewDecoder(req.Body)
	var p param
	err = decoder.Decode(&p)
	if err != nil {
		respondErrorWriter(w, http.StatusInternalServerError, "Error trying to decode de request", err)
		return
	}

	NewHash, err := auth.HashPassword(p.Password)
	if err != nil {
		respondErrorWriter(w, http.StatusInternalServerError, "Error hashing new password", err)
		return
	}
	updatedUser, err := a.Queries.UpdateUser(req.Context(), database.UpdateUserParams{
		HashedPassword: NewHash,
		Email:          p.Email,
		ID:             user})

	responWithJson(w, http.StatusOK, User{
		Id:            updatedUser.ID,
		Created_at:    updatedUser.CreatedAt,
		Updated_at:    updatedUser.UpdatedAt,
		Email:         updatedUser.Email,
		Is_Chirpy_Red: updatedUser.IsChirpyRed.Bool,
	})
}
