package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func MakeJWT(user uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject:   user.String(),
	})
	tokenSigned, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}
	return tokenSigned, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	//Pasas un pointer a los Registered que pasaste arriba para que los llenes con los claims que pasaste
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.Nil, err
	}
	userID, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, err
	}
	ID, err := uuid.Parse(userID)
	if err != nil {
		return uuid.Nil, err
	}
	return ID, nil
}
