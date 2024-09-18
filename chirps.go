package main

import (
	"encoding/json"
	"net/http"
)

type parameters struct {
	Body string `json:"body"`
}

type returnErr struct {
	Error string `json:"error"` 
}

type okStruct struct {
	Valid string `json:"cleaned_body"`
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

func (c *apiConfig) createChripHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	switch r.Method {
	case http.MethodGet:
		// Handle Get request 
		chirps, err := c.DB.GetChirps()
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "failed to retrieve chirp")
			return 
		} 
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(chirps)
		
	case http.MethodPost:
		var params parameters
		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
				respondWithError(w, http.StatusBadRequest, "invalid request payload")
				return 
			}
			if len(params.Body) == 0 {
				respondWithError(w, http.StatusBadRequest, "chirp cannot be empty")
				return
			}
			if len(params.Body) > 140 {
				respondWithError(w, http.StatusBadRequest, "chirp cannot be longer than 140 characters")
				return
			}
			chirp, err := c.DB.CreateChirp(params.Body)
			
			if err != nil {
				respondWithError(w, http.StatusInternalServerError, "failed to create chirp")
				return 
			}
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(chirp)
			
	default :
			respondWithError(w, http.StatusMethodNotAllowed, "method not supported")
			return
		}

	
}