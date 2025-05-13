package auth

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

func CheckPasswordHash(hash, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	return err
}