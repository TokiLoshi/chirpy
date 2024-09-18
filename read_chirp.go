package main

import (
	"fmt"
	"net/http"
)

func (c *apiConfig) readChirp(w http.ResponseWriter, r *http.Request) {
	fmt.Println("getting tweets")
}