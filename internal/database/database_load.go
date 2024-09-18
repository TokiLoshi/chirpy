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
	
	// Read the file contents
	data, err := os.ReadFile(db.path)
	if err != nil {
		return loadedDB, fmt.Errorf("error opening db, %w", err)
	}

	if len(data) == 0 {
		fmt.Printf("database is empty")
	}

	err = json.Unmarshal(data, &loadedDB)
	if err != nil {
		return  loadedDB, fmt.Errorf("error marshalling everything in load: %w", err)
	}
	return loadedDB, nil
}