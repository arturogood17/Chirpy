package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type TokenType string

const (
	Tokentype TokenType = "chirpy-access"
)

func HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "error hashing the password", err
	}
	return string(hashed), nil
}

func CheckPasswordHash(hash, password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return fmt.Errorf("password incorrect: %w", err)
	}
	return nil
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    string(Tokentype),
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject:   userID.String(),
	})
	return token.SignedString([]byte(tokenSecret))
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claims := jwt.RegisteredClaims{}
	tokenVal, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.Nil, err
	}
	ID, err := tokenVal.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, err
	}
	issuer, err := tokenVal.Claims.GetIssuer()
	if err != nil {
		return uuid.Nil, err
	}
	if issuer != string(Tokentype) {
		return uuid.Nil, errors.New("invalid issuer")
	}
	id, err := uuid.Parse(ID)
	if err != err {
		return uuid.Nil, err
	}
	return id, nil
}

func MakeRefreshToken() (string, error) {
	token := make([]byte, 32)
	rand.Read(token)
	return hex.EncodeToString(token), nil
}

func GetBearerToken(headers http.Header) (string, error) {
	auth := strings.TrimPrefix(headers.Get("Authorization"), "Bearer ")
	if auth == "" {
		return "", errors.New("no authentication token was provided")
	}
	return auth, nil
}
