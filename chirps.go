package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	auth "github.com/TokiLoshi/chirpy/internal"
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

	bearerToken, err := auth.GetBearerToken(req.Header)
	
	if err != nil {
		respondWithError(w, 401, "Unauthorized")
		return
	}

	validatedUserId, err := auth.ValidateJWT(bearerToken, cfg.jwtSecret)
	
	if err != nil {
		respondWithError(w, 401, "Unauthorized")
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

	
	newChirp, err := cfg.dbQueries.CreateChirp(req.Context(), database.CreateChirpParams{
		Body: cleanchirp,
		UserID: validatedUserId,
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

type ChirpResponse struct {
	ID uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body string `json:"body"`
	UserId uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) getAllChirps(w http.ResponseWriter, req *http.Request) {
	allChirps, err := cfg.dbQueries.GetAllChirps(req.Context())
	if err != nil {
		respondWithError(w, 400, "couldn't getAllChirps" + err.Error())
		return
	}

	responseChirps := make([]ChirpResponse, len(allChirps))
	for i, chirp := range allChirps {
		responseChirps[i] = ChirpResponse {
			ID: chirp.ID, 
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body: chirp.Body,
			UserId: chirp.UserID,
		}
	}


	respondWithJson(w, 200, responseChirps)


}

func (cfg *apiConfig) getSingleChirp(w http.ResponseWriter, req *http.Request) {
	chirpId := req.PathValue("chirpID")
	fmt.Printf("path: %+v\n", chirpId)
	if len(chirpId) == 0 {
		respondWithError(w, 400, "invalid chirp ")
		return
	}

	id, err := uuid.Parse(chirpId)
	if err != nil {
		respondWithError(w, 400, "invalid id")
		return
	}

	singleChirp, err := cfg.dbQueries.GetSingleChirp(req.Context(), id)
	if err != nil {
		respondWithError(w, 404, "couldn't get chirp" + err.Error())
		return
	}


	responseChirp := ChirpResponse {
		ID: singleChirp.ID,
		CreatedAt: singleChirp.CreatedAt,
		UpdatedAt: singleChirp.UpdatedAt, 
		Body: singleChirp.Body, 
		UserId: singleChirp.UserID,
	}

	respondWithJson(w, 200, responseChirp)
	

}

func (cfg *apiConfig) handleDelete(w http.ResponseWriter, req *http.Request) {

	chirpId := req.PathValue("chirpID")
	if len(chirpId) == 0 {
		respondWithError(w, 404, "Invalid Id")
		return
	}
	// This is authentcated so check token in header
	authHeader := req.Header.Get("Authorization")
	if authHeader == "" {
		respondWithError(w, 401, "Auth header is missing")
		return
	}

	// Format for Bearer
	fields := strings.Split(authHeader, " ")
	if len(fields) != 2 || fields[0] != "Bearer" {
		respondWithError(w, 401, "Malformed Authorization")
		return
	}

	tokenString := fields[1] 

	userId, err := auth.ValidateJWT(tokenString, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, 401, "Invalid token")
		return
	}

	id, err := uuid.Parse(chirpId)

	chirpInfo, err := cfg.dbQueries.GetSingleChirp(req.Context(), id)
	if err != nil {
		respondWithError(w, 404, "Chirp does not exist")
		return
	}

	if chirpInfo.UserID != userId {
		respondWithError(w, 403, "Unauthorized")
		return
	}

	_, err = cfg.dbQueries.DeleteChirp(req.Context(), id)
	if err != nil {
		respondWithError(w, 404, "chirp could not be found")
		return
	}

	w.WriteHeader(204)

}
