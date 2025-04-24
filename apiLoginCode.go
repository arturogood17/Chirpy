package main

import (
	"encoding/json"
	"net/http"

	"github.com/arturogood17/Chirpy/internal/auth"
)

func (a *apiConfig) hLogin(w http.ResponseWriter, req *http.Request) {
	type Param struct {
		Password string `json:"password"`
		Email    string `json:"email"`
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
	responWithJson(w, http.StatusOK, User{Id: user.ID,
		Created_at: user.CreatedAt,
		Updated_at: user.UpdatedAt,
		Email:      user.Email})
}
