package database

import (
	"encoding/json"
	"fmt"
	"os"
)

// load db reads database file into memory
func (db *DB) LoadDB() (DBStructure, error) {
	
	fmt.Println("reading db file into memory")
	
	var loadedDB DBStructure
	loadedDB.Chirps = make(map[int]Chirp)
	
	// Read the file contents
	data, err := os.ReadFile(db.path)
	if err != nil {
		return loadedDB, fmt.Errorf("error opening db, %w", err)
	}
	fmt.Printf("File contents: %s\n", data)
	if len(data) == 0 {
		fmt.Printf("database is empty")
		return loadedDB, nil
	}

	err = json.Unmarshal(data, &loadedDB)
	if err != nil {
		fmt.Printf("error in un marshal: %v", err)
		return  loadedDB, fmt.Errorf("error marshalling everything in load: %w", err)
	}
	return loadedDB, nil
}