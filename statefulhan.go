package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/arturogood17/Chirpy/internal/auth"
	"github.com/arturogood17/Chirpy/internal/database"
	"github.com/google/uuid"
)

func (a *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a.hits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (a *apiConfig) HandlerMetrics(res http.ResponseWriter, req *http.Request) {
	req.Header.Set("Content-Type", "text/html")
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(fmt.Sprintf(`<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>
`, a.hits.Load())))
}

func (a *apiConfig) HandlerReset(res http.ResponseWriter, req *http.Request) {
	if a.PLATFORM != "dev" {
		log.Fatal("403 Forbidden")
	}
	if err := a.dbQueries.DeleteUsers(req.Context()); err != nil {
		log.Fatalf("error deleting the users from the database")
	}
}

func (a *apiConfig) HandlerUser(res http.ResponseWriter, req *http.Request) {
	if req.Body == nil {
		respondWithError(res, http.StatusBadRequest, "Empty request", nil)
	}
	type UserCreation struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	var UserData UserCreation
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&UserData); err != nil {
		respondWithError(res, 500, "Error decoding Json", err)
		return
	}
	if UserData.Password == "" {
		respondWithError(res, 400, "User needs a password", nil)
		return
	}
	hashedPass, err := auth.HashPassword(UserData.Password)
	if err != nil {
		respondWithError(res, 500, "Error hashing the password", err)
		return
	}
	user, err := a.dbQueries.CreateUser(req.Context(), database.CreateUserParams{
		Email:          UserData.Email,
		HashedPassword: hashedPass,
	})
	if err != nil {
		respondWithError(res, 500, "Error creating user", err)
		return
	}
	new_user := User{
		ID:          user.ID.String(),
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed.Bool,
	}
	respondWithJson(res, 201, new_user)
}

func (a *apiConfig) HandlerChirps(res http.ResponseWriter, req *http.Request) {
	type reqChirp struct {
		Body string `json:"body"`
	}
	TokenAuth, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(res, 401, "Coukldn't find JWT", err)
		return
	}
	userID, err := auth.ValidateJWT(TokenAuth, a.SECRET)
	if err != nil {
		respondWithError(res, 401, "O este", err)
		return
	}

	decoder := json.NewDecoder(req.Body)
	chirp := reqChirp{}
	if err := decoder.Decode(&chirp); err != nil {
		respondWithError(res, 500, "Error decoding JSON", err)
		return
	}
	if len(chirp.Body) > 140 || len(chirp.Body) == 0 {
		respondWithError(res, http.StatusBadGateway, "Chirp is too long or is empty", nil)
		return
	}
	bodyVal := WordValidation(chirp.Body)

	new_chirp, err := a.dbQueries.CreateChirps(req.Context(), database.CreateChirpsParams{
		Body:   bodyVal,
		UserID: userID,
	})

	if err != nil {
		respondWithError(res, 500, "Couldn't create the chirp", err)
		return
	}
	nc := Chirp{
		ID:        new_chirp.ID.String(),
		CreatedAt: new_chirp.CreatedAt,
		UpdatedAt: new_chirp.UpdatedAt,
		Body:      new_chirp.Body,
		UserID:    new_chirp.UserID.String(),
	}
	respondWithJson(res, 201, nc)
}

func (a *apiConfig) AllChirps(res http.ResponseWriter, req *http.Request) {
	author_id := req.URL.Query().Get("author_id")
	if author_id != "" {
		aID, err := uuid.Parse(author_id)
		if err != nil {
			respondWithError(res, 500, "Error parsing author_id", err)
			return
		}
		chirps, err := a.dbQueries.ChirpByAuthor(req.Context(), aID)
		if err != nil {
			respondWithError(res, 404, "Error getting chirps", err)
		}
		SChirpsByAuthor := MappingChirps(chirps)
		respondWithJson(res, 200, SChirpsByAuthor)
	} else {
		chirps, err := a.dbQueries.AllChirps(req.Context())
		if err != nil {
			respondWithError(res, 500, "Error retrieving the all chirps from database", err)
			return
		}
		SChirps := MappingChirps(chirps)
		respondWithJson(res, 200, SChirps)
	}
}

func MappingChirps(chirps []database.Chirp) []Chirp {
	var slice_chirps []Chirp
	for _, chirp := range chirps {
		nc := Chirp{
			ID:        chirp.ID.String(),
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID.String(),
		}
		slice_chirps = append(slice_chirps, nc)
	}
	return slice_chirps
}

func (a *apiConfig) SingleChirp(res http.ResponseWriter, req *http.Request) {
	pathValue := req.PathValue("chirpID")
	if pathValue == "" {
		respondWithError(res, 404, "Chirp not found", nil)
		return
	}
	chirpID, err := uuid.Parse(pathValue)
	if err != nil {
		respondWithError(res, 500, "Error parsing the chirpID", err)
		return
	}

	chirp, err := a.dbQueries.SingleChirp(req.Context(), chirpID)
	if err != nil {
		respondWithError(res, 404, "Error fetching the desired chirp", err)
		return
	}
	nc := Chirp{
		ID:        chirp.ID.String(),
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID.String(),
	}
	respondWithJson(res, 200, nc)
}

func (a *apiConfig) UserLogin(res http.ResponseWriter, req *http.Request) {
	if req.Body == nil {
		respondWithError(res, http.StatusBadRequest, "Empty request", nil)
	}
	type LoginData struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	type response struct {
		User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}
	var UserData LoginData
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&UserData); err != nil {
		respondWithError(res, 500, "Error decoding Json", err)
		return
	}
	user, err := a.dbQueries.SearchUser(req.Context(), UserData.Email)
	if err != nil {
		respondWithError(res, 401, "Incorrect email or password", err)
		return
	}
	if err = auth.CheckPasswordHash(user.HashedPassword, UserData.Password); err != nil {
		respondWithError(res, 401, "Incorrect email or password", err)
		return
	}
	token, err := auth.MakeJWT(user.ID, a.SECRET, time.Hour)
	if err != nil {
		respondWithError(res, 500, "Error creating token", err)
		return
	}

	Rtoken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(res, 500, "Couldn't create refresh token", err)
		return
	}

	_, err = a.dbQueries.CreateRefreshToken(req.Context(), database.CreateRefreshTokenParams{
		Token:     Rtoken,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(1440 * time.Hour),
	})
	if err != nil {
		respondWithError(res, 500, "Couldn't save refresh token in database", err)
		return
	}

	r := response{
		User: User{
			ID:          user.ID.String(),
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
			Email:       user.Email,
			IsChirpyRed: user.IsChirpyRed.Bool,
		},
		Token:        token,
		RefreshToken: Rtoken,
	}
	respondWithJson(res, 200, r)
}

func (a *apiConfig) RefreshToken(res http.ResponseWriter, req *http.Request) {
	type OneHToken struct {
		Token string `json:"token"`
	}
	refreshToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(res, 500, "No refresh token found in request", err)
		return
	}
	user, err := a.dbQueries.GetUserFromRefreshToken(req.Context(), refreshToken)
	if err != nil {
		respondWithError(res, 401, "Error getting user from refresh token", err)
		return
	}
	newAtoken, err := auth.MakeJWT(user.ID, a.SECRET, time.Hour)
	if err != nil {
		respondWithError(res, 500, "Error creating a new access token", err)
		return
	}
	respondWithJson(res, 200, OneHToken{
		Token: newAtoken,
	})

}

func (a *apiConfig) Revoke(res http.ResponseWriter, req *http.Request) {
	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(res, 500, "No refresh token found in request", err)
		return
	}
	if err = a.dbQueries.Revoke(req.Context(), token); err != nil {
		respondWithError(res, 500, "Error revoking token", err)
		return
	}
	res.WriteHeader(http.StatusNoContent)
}

func (a *apiConfig) UpdateUser(res http.ResponseWriter, req *http.Request) {
	if req.Body == nil {
		respondWithError(res, 401, "Body of request missing", nil)
		return
	}
	type UpdatedInfo struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	var UpInfo UpdatedInfo
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&UpInfo); err != nil {
		respondWithError(res, 500, "Error decoding body of request", err)
	}
	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(res, 401, "no token found", err)
		return
	}
	user, err := auth.ValidateJWT(token, a.SECRET)
	if err != nil {
		respondWithError(res, 401, "Invalid token", err)
		return
	}
	hashedP, err := auth.HashPassword(UpInfo.Password)
	if err != nil {
		respondWithError(res, 500, "Error hashing password", err)
		return
	}
	updatedU, err := a.dbQueries.UpdateUser(req.Context(), database.UpdateUserParams{Email: UpInfo.Email,
		HashedPassword: hashedP,
		ID:             user,
	})
	if err != nil {
		respondWithError(res, 500, "Error updating user", err)
		return
	}
	respondWithJson(res, 200, User{
		ID:          updatedU.ID.String(),
		CreatedAt:   updatedU.CreatedAt,
		UpdatedAt:   updatedU.UpdatedAt,
		Email:       UpInfo.Email,
		IsChirpyRed: updatedU.IsChirpyRed.Bool,
	})
}

func (a *apiConfig) DeleteChirp(res http.ResponseWriter, req *http.Request) {
	pathValue := req.PathValue("chirpID")
	if pathValue == "" {
		respondWithError(res, 404, "Chirp not found", nil)
		return
	}
	chirpID, err := uuid.Parse(pathValue)
	if err != nil {
		respondWithError(res, 500, "Couldn't parse chirpID", err)
		return
	}
	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(res, 401, "Token not found", err)
		return
	}
	ValUser, err := auth.ValidateJWT(token, a.SECRET)
	if err != nil {
		respondWithError(res, 403, "Invalid token", err)
		return
	}
	chirp, err := a.dbQueries.SingleChirp(req.Context(), chirpID)
	if err != nil {
		respondWithError(res, 404, "Chirp not found", err)
		return
	}
	if chirp.UserID != ValUser {
		respondWithError(res, 403, "You can't delete this chirp", nil)
		return
	}
	if err = a.dbQueries.DeleteChirp(req.Context(), chirpID); err != nil {
		respondWithError(res, 500, "Couldn't delete chirp", err)
		return
	}
	res.WriteHeader(204)
}

func (a *apiConfig) RedChirpy(res http.ResponseWriter, req *http.Request) {
	polkaK, err := auth.GetAPIKey(req.Header)
	if err != nil {
		respondWithError(res, 401, "authentication key for polka was not found", err)
		return
	}
	if polkaK != a.POLKAKEY {
		respondWithError(res, 401, "Not authorized", nil)
		return
	}
	type RedChirpyRequest struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}
	var red RedChirpyRequest
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&red); err != nil {
		respondWithError(res, 500, "Error decoding body of request", err)
	}
	if red.Event != "user.upgraded" {
		res.WriteHeader(204)
		return
	}
	userID, err := uuid.Parse(red.Data.UserID)
	if err != nil {
		respondWithError(res, 500, "error parsing userID", err)
		return
	}
	if err = a.dbQueries.WelcomeToChirpy(req.Context(), userID); err != nil {
		respondWithError(res, 404, "User not found", err)
		return
	}
	respondWithJson(res, 204, nil)
}
