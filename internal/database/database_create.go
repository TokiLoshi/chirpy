package database

import (
	"fmt"
)

// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateChirp(body string) (Chirp, error) {
	fmt.Println("Creating a new chirp")
	
	newChirp := Chirp{}
	// Load the database
	data, err := LoadDB()
	if err != nil {
		return newChirp, fmt.Errorf("error loading db: %w", err)
	}
	fmt.Println(data)

	// Validate the chirp 

	// Write the chirp 
	return newChirp, nil
}
