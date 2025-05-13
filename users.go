package main

import (
	"encoding/json"
	"net/http"
	"time"

	auth "github.com/TokiLoshi/chirpy/internal"
	"github.com/TokiLoshi/chirpy/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email string `json:"email"`
}

func (cfg *apiConfig) validateUser(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}

	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, "error decoding json data")
		return
	}

	userEmail := params.Email
	userPassword := params.Password 

	hashedPassword, err := auth.HashPassword(userPassword)

	if err != nil {
		respondWithError(w, 500, "unable to hash password")
	}

	dbUser, err := cfg.dbQueries.CreateUser(req.Context(), database.CreateUserParams{
		Email: userEmail, 
		HashedPassword: hashedPassword,
	})

	if err != nil {
		respondWithError(w, 400, "couldn't create user: " + err.Error() )
		return
		} 
	
	user := User {
		ID: dbUser.ID,
		Email: dbUser.Email,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
	}

	respondWithJson(w, 201, user)

}

func (cfg *apiConfig) handleLogin(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}
	decoder := json.NewDecoder(req.Body)
	params := parameters{}

	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, "error decoding json data")
		return
	}

	userEmail := params.Email 
	userPassword := params.Password 

	userEntry, err := cfg.dbQueries.GetUserByEmail(req.Context(), userEmail)

	if err != nil {
		respondWithError(w, 401, "Unauthorized")
		return
	}

	error := auth.CheckPasswordHash(userEntry.HashedPassword, userPassword)

	if error != nil {
		respondWithError(w, 401, "Unauthorized")
		return
	}

	user := User{
		ID: userEntry.ID,
		Email: userEntry.Email, 
		CreatedAt: userEntry.CreatedAt, 
		UpdatedAt: userEntry.UpdatedAt,

	}

	respondWithJson(w, 200, user)

}