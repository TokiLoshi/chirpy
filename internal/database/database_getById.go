package database

import (
	"fmt"
)

func (db *DB) GetChirpsById(id int) (Chirp, error) {
	fmt.Printf("Looking for a chirp with id: %v\n", id)
	db.mux.Lock()
	defer db.mux.Unlock()
	
	specificChirp := Chirp{}

	dbStructure, err := db.LoadDB()
	if err != nil {
		return specificChirp, fmt.Errorf("couldn't load db: %w", err)
	}
	chirp, exists := dbStructure.Chirps[id] 
	if !exists {
		return specificChirp, fmt.Errorf("chirp with id: %v does not exist", id)
	}
	return chirp, nil 
}