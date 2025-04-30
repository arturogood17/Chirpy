package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/arturogood17/Chirpy/internal/auth"
	"github.com/arturogood17/Chirpy/internal/database"
)

func (a *apiConfig) hLogin(w http.ResponseWriter, req *http.Request) {
	type Param struct {
		Password         string        `json:"password"`
		Email            string        `json:"email"`
		ExpiresInSeconds time.Duration `json:"expires_in_seconds"`
	}
	decoder := json.NewDecoder(req.Body)
	var param Param
	err := decoder.Decode(&param)
	if err != nil {
		respondErrorWriter(w, http.StatusInternalServerError, "Error trying to decode de request", err)
		return
	}
	if param.Email == "" || param.Password == "" {
		respondErrorWriter(w, http.StatusBadRequest, "None of the fileds should be empty", nil)
		return
	}
	user, err := a.Queries.GetUserByEmail(req.Context(), param.Email)
	if err != nil {
		respondErrorWriter(w, http.StatusInternalServerError, "Error getting user to login", err)
		return
	}
	err = auth.CheckPasswordHash(user.HashedPassword, param.Password)
	if err != nil {
		respondErrorWriter(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}
	token, err := auth.MakeJWT(user.ID, a.SECRET)
	if err != nil {
		respondErrorWriter(w, http.StatusInternalServerError, "Error creating JWT token", err)
		return
	}
	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondErrorWriter(w, http.StatusInternalServerError, "Error creating refresh token", err)
		return
	}
	RToken, err := a.Queries.CreateRefreshToken(req.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    user.ID,
		ExpiresAt: time.Now().UTC().Add((time.Hour * 24) * 60),
	})
	if err != nil {
		respondErrorWriter(w, http.StatusInternalServerError, "Error saving refresh token", err)
		return
	}
	type response struct {
		User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}
	responWithJson(w, http.StatusOK, response{
		User: User{Id: user.ID,
			Created_at: user.CreatedAt,
			Updated_at: user.UpdatedAt,
			Email:      user.Email},
		Token:        token,
		RefreshToken: RToken.Token,
	})
}
