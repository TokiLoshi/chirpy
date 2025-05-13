package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
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



func MakeJWT(userId uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)), 
		Issuer: "chirpy",
		IssuedAt: jwt.NewNumericDate(time.Now()),
		Subject: userId.String(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "could not sign token", err
	}

	return signedToken, nil

}


func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	type UserClaims struct {
		ExpiresAt time.Time;
		Issuer string;
		IssuedAt time.Time;
		Subject string;
	}

	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})

	if err != nil {
		return uuid.Nil, err
	} 
	
	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		userId, err := uuid.Parse(claims.Subject)
		if err != nil {
			return uuid.Nil, err
		}
		return userId, nil
	}
	return uuid.Nil, fmt.Errorf("Invalid token claims")
}

func GetBearerToken(headers http.Header) (string, error) {
	// Auth information coms iin through the Authorization Header 
	// Value will be Bearer TOKEN_STRING 
	authInfo := headers.Get("Authorization")
	if authInfo == "" {
		return "", fmt.Errorf("No authorization")
	}
	// Check for Authorization in the headers and return Token String
	fields := strings.Fields(authInfo)
	if len(fields) < 2 || strings.ToLower(fields[0]) != "bearer" {
		return "", fmt.Errorf("No bearer token")
	} 
	tokenString := fields[1]
	return tokenString, nil 
	// write a unit thest for it 

}

func MakeRefreshToken() (string, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key) 

	if err != nil {
		newError := fmt.Errorf("error generating random noise for refresh token: %v", err)

		return "", newError
	}
	stringifiedToken := hex.EncodeToString(key)

	return stringifiedToken, nil
}
