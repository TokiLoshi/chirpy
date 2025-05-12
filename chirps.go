package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/TokiLoshi/chirpy/internal/database"
	"github.com/google/uuid"
)

type Chirp struct {
	ID uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body string `json:"body"`
	UserId uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) createChirp(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Body string `json:"body"`
		User_Id string `json:"user_id"`
	}

	// Decode the request 
	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	
	// Handle errors decoding the request 
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, "error decoding json data")
		return 
	}

	// Check if chirps length is too long 
	maxLength := 140
	if len(params.Body) > maxLength {
		respondWithError(w, 400, "Chirp is too long")
		return 
	}

	// Check if the chirp is clean 
	cleanchirp := cleanChirpBody(params.Body)

	userUUID, err := uuid.Parse(params.User_Id)
	if err != nil {
		respondWithError(w, 400, "invalid user_id")
	}
	
	newChirp, err := cfg.dbQueries.CreateChirp(req.Context(), database.CreateChirpParams{
		Body: cleanchirp,
		UserID: userUUID,
	})

	if err != nil {
		respondWithError(w, 400, "couldn't create chirp: " + err.Error())
		return
	}

	chirp := Chirp {
		ID: newChirp.ID,
		CreatedAt: newChirp.CreatedAt,
		UpdatedAt: newChirp.UpdatedAt,
		Body: newChirp.Body,
		UserId: newChirp.UserID,
	}

	respondWithJson(w, 201, chirp)

}

// type cleanedResponse struct {
// 	CleanedBody string `json:"cleaned_body"`
// }

func cleanChirpBody(chirp string) string {
	
	
	words := strings.Split(chirp, " ")
	cleanedWords := make([]string, 0, len(words))
	forbiddenWords := map[string]bool{
		"kerfuffle": true, 
		"sharbert": true, 
		"fornax": true,
	}
	for _, word := range words {
		loweredWord := strings.ToLower(word)
		if forbiddenWords[loweredWord] {
			cleanedWords = append(cleanedWords, "****")
		} else {
			cleanedWords = append(cleanedWords, word)
		}
		
	}
	cleanChirp := strings.Join(cleanedWords, " ")
	return cleanChirp
}

