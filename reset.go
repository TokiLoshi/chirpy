package main

import (
	"fmt"
	"log"
	"net/http"
)

func (c *apiConfig) reset() {
	c.fileserveHits = 0
}

func (c *apiConfig) resetHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("resetting fileserveHits")
	c.reset()
	w.Header().Set("Content-Type", "text/lain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Status reset to %d", c.fileserveHits)
}