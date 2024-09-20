package main

import (
	"encoding/json"
	"net/http"
	"strconv"
)


func (c *apiConfig) singleChirpHandler(w http.ResponseWriter, r * http.Request) {
	chirpId := r.PathValue("chirpId")

	if len(chirpId) == 0 {
		respondWithError(w, http.StatusInternalServerError, "id cannot be empty") 
	}
	intChirpId, err := strconv.Atoi(chirpId)
	
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "something is wrong with the id")
	}

	chirp, err := c.DB.GetChirpsById(intChirpId)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "chirp doesn't exist")
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(chirp)
	
	
}