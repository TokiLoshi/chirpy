package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)


func (c *apiConfig) singleChirpHandler(w http.ResponseWriter, r * http.Request) {
	fmt.Println("looking for a single chirp")
	chirpId := r.PathValue("chirpId")
	fmt.Printf("Path value returned: %v\n", chirpId)
	if len(chirpId) == 0 {
		respondWithError(w, http.StatusInternalServerError, "id cannot be empty") 
	}
	intChirpId, err := strconv.Atoi(chirpId)
	
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "something is wrong with the id")
	}
	fmt.Printf("We have a chirp to look for: %v\n", intChirpId)

	chirp, err := c.DB.GetChirpsById(intChirpId)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "chirp doesn't exist")
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(chirp)
	
	
}