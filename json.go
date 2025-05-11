package main

import (
	"encoding/json"
	"log"
	"net/http"
)
func respondWithError(w http.ResponseWriter, code int, message string) {
	type errorResponse struct {
		Error string `json:"error"`
	}
	responseError := errorResponse{
		Error: message,
	}
	data, err := json.Marshal(responseError)
	if err != nil {
		log.Printf("error marshalling JSON %s", err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	dat, err := json.Marshal(payload)
	if err != nil {
		respondWithError(w, 400, "error marshalling JSON")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)	
	w.Write(dat)
}