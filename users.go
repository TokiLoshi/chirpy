package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type userParameters struct {
	Body string `json:"email"`
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

	if len(params.Body) == 0 {
		respondWithError(w, http.StatusBadRequest, "email cannot be empty")
		return
	}

	email, err := c.DB.CreateUser(params.Body)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error creating user")
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(email)

}