package database

import (
	"fmt"
	"sort"
)

// Returns all Chirps in the database
func (db *DB) GetChirps() ([]Chirp, error) {
	fmt.Println("Getting all the chirps")
	db.mux.Lock()
	defer db.mux.Unlock()
	
	allChirps := []Chirp{}
	dbStructure, err := db.LoadDB()
	if err != nil {
		return allChirps, fmt.Errorf("error loading chirps: %w", err)
	}

	for _, chirp := range dbStructure.Chirps {
		allChirps = append(allChirps, chirp)
	}

	sort.Slice(allChirps, func(i, j int) bool {
		return allChirps[i].Id < allChirps[j].Id
	})

	return allChirps, nil
}
