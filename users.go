package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type userParameters struct {
	Password string `json:"password"`
	Email string `json:"email"`
	Expires string `json:"expires_in_seconds"`
}

type UserResponse struct {
	Id int `json:"id"`
	Email string `json:"email"`
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash),[]byte(password))
	return err == nil
}

	
	// Set Subject to stringified version of user id
	func (c *apiConfig) CreateToken(userId string, expiry int) (string, error) {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, 
		jwt.MapClaims {
			"issuer" : "chirpy",
			"iat": time.Now().Unix(),
			// Set ExpiresAt to the current time plus expiration (expires_in_seconds)
			"exp" : time.Now().Add(time.Duration(expiry) * time.Second).Unix(),
			"subject" : string(userId),
		})
		tokenString, err := token.SignedString([]byte(c.jwt))
		if err != nil {
			return "", fmt.Errorf("error signing token string from env: %w", err)
		}
		return tokenString, nil
	}

func (c *apiConfig) userLoginHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("User trying to login")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	var params userParameters 
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request")
		return
	}

	email := params.Email
	password := params.Password 
	expiry := params.Expires
	fmt.Printf("expiry token from params: %v\n", expiry)

	fmt.Printf("params: %v\n", params)
	if len(email) == 0 {
		respondWithError(w, http.StatusBadRequest, "email cannot be empty")
		return
	}

	if len(password) == 0 {
		respondWithError(w, http.StatusBadRequest, "password cannot be empty")
		return 
	}

	// Load Database 
	// Check to see if email exists 
	user, err := c.DB.GetUser(email)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error checking on user")
		return
	}
	if user == nil {
		fmt.Printf("User doesn't exist: %v\n", user)
		respondWithError(w, http.StatusUnauthorized, "user doesn't exist")
		return
	}

	
	fmt.Printf("Here's the user stuff: %v\n", user)

	hashedPassword := user.Password
	matches := CheckPasswordHash(password, hashedPassword)

	
	if !matches {
		fmt.Printf("incorrect password: %v\n", matches)
		respondWithError(w, http.StatusUnauthorized, "incorrect password")
	} 

	testToken, err := c.CreateToken(string(user.Id), 2)
	if err != nil {
		fmt.Printf("error creating test token: %v\n", err)
	}
	fmt.Printf("Test token: %v\n", testToken)


	// 2. Use token 
	// Use token.SignedString to sign the token with the secret key 
	// If expires_in_seconds is not specified  make the default 24 hours 
	// If client specified number over 24 hours use 24 hours 
	// Add token to user response 
	
	response := UserResponse {
		Id : user.Id,
		Email : user.Email,
	} 
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)


}


func (c *apiConfig) userHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Trying to get users")
	// Allows users to be created 
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	var params userParameters 
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request payload")
		return
	}

	fmt.Printf("params: %v\n", params)

	if len(params.Email) == 0 {
		respondWithError(w, http.StatusBadRequest, "email cannot be empty")
		return
	}

	if len(params.Password) == 0 {
		respondWithError(w, http.StatusBadRequest, "email cannot be empty")
		return
	}


	hashedPassword, err := HashPassword(params.Password)
	if err != nil {
		fmt.Printf("Error hashing password: %v\n", hashedPassword)
	}

	newUser, err := c.DB.CreateUser(params.Email, hashedPassword)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error creating user")
		return
	}
	
	response := UserResponse {
		Id : newUser.Id,
		Email : newUser.Email,
	} 

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)

}