package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func validateChirps(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Body string `json:"body"`
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
	respondWithJson(w, 200, cleanedResponse{cleanchirp})

}


type cleanedResponse struct {
	CleanedBody string `json:"cleaned_body"`
}

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

