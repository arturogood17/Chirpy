package main

import (
	"net/http"

	"github.com/arturogood17/Chirpy/internal/auth"
)

func (a *apiConfig) RefreshToken(w http.ResponseWriter, req *http.Request) {
	type ACCToken struct {
		Token string `json:"token"`
	}
	authRtoken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondErrorWriter(w, http.StatusBadRequest, "error with refresh token", err)
		return
	}
	user, err := a.Queries.GetUserFromRefreshToken(req.Context(), authRtoken)
	if err != nil {
		respondErrorWriter(w, http.StatusUnauthorized, "error getting user with refresh token", err)
		return
	}

	accToken, err := auth.MakeJWT(user.ID, a.SECRET)
	if err != nil {
		respondErrorWriter(w, http.StatusInternalServerError, "error creating access token for user", err)
		return
	}
	responWithJson(w, http.StatusOK, ACCToken{Token: accToken})
}

func (a *apiConfig) RefreshRevoke(w http.ResponseWriter, req *http.Request) {
	authRtoken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondErrorWriter(w, http.StatusBadRequest, "error with refresh token", err)
		return
	}
	err = a.Queries.SetRevoke(req.Context(), authRtoken)
	if err != nil {
		respondErrorWriter(w, http.StatusInternalServerError, "error revoking token", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
