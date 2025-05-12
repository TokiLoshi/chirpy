package main

import (
	"net/http"
	"os"
)

func (cfg *apiConfig) adminReset(w http.ResponseWriter, req *http.Request) {
	platformType := os.Getenv("PLATFORM")
	if platformType != "dev" {
		respondWithError(w, 403, "Forbidden")
		return
	}

	err := cfg.dbQueries.DeleteUsers(req.Context())
	if err != nil {
		respondWithError(w, 500, "couldn't delete users")
		return
	}

	respondWithJson(w, 200, struct{}{})

}