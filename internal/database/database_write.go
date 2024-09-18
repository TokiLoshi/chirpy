package database

import (
	"encoding/json"
	"fmt"
	"os"
)

// Writes database file to disk
func (db *DB) WriteDB(dbStructure DBStructure) error {
	fmt.Printf("We're getting into the write db with: %v\n", dbStructure)
	// db.mux.Lock()
	// defer db.mux.Unlock()

	fmt.Printf("writing file to disk")
	filepath := db.path 

	// Marshall the data 
	newJson, err := json.Marshal(dbStructure)
	if err != nil {
		return fmt.Errorf("error marshalling file: %w", err)
	} 

	err = os.WriteFile(filepath, newJson, 0666)
	if err != nil {
		return fmt.Errorf("error writing file: %w", err)
	}
	
	
	
	return nil
}