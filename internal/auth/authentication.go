package auth

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) { //convierte el password y te da un hash para almacenarlo
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func CheckPasswordHash(hashedPassword, password string) error { //compara el hash con el password almacenado a ver si coinciden
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return fmt.Errorf("Invalid user - %v", err)
	}
	return nil
}
