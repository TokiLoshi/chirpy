package database

import (
	"fmt"
)

// Returns all Chirps in the database
func (db *DB) GetChirps() ([]Chirp, error) {
	fmt.Println("Getting all the chirps")
	allChirps := []Chirp{}
	return allChirps, nil
}