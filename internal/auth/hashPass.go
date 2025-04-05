package auth

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type TokenType string

const (
	tokentype TokenType = "chirpy-access"
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
		Issuer:    "chirpy-access",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject:   userID.String(),
	})
	return token.SignedString([]byte(tokenSecret))
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	var userID uuid.UUID
	claims := jwt.RegisteredClaims{}
	tokenVal, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return userID, err
	}
	ID, err := tokenVal.Claims.GetSubject()
	if err != nil {
		return userID, err
	}
	issuer, err := tokenVal.Claims.GetIssuer()
	if err != nil {
		return userID, err
	}
	if issuer != string(tokentype) {
		return userID, errors.New("invalid issuer")
	}
	id, err := uuid.Parse(ID)
	if issuer != string(tokentype) {
		return userID, err
	}
	return id, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	headers.Get("Authorization")
}
