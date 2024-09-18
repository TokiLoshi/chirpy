package database

import (
	"encoding/json"
	"fmt"
	"os"
)

// ensureDB creates new databse file if it doesn't exist
func (db *DB) ensureDB() error {
	fmt.Println("creating the databse file")
	filepath := db.path

	// If there's an error we want to create a file
	_, err := os.Stat(db.path)
	if os.IsNotExist(err) {
		fmt.Printf("file doesn't exist, writing a new one: %v\n", err)
		newDB := &DBStructure{}
		newDB.Chirps = make(map[int]Chirp)
		newJson, err := json.Marshal(newDB)
		if err != nil {
			return fmt.Errorf("error marshalling new db: %w", err)
		}
		err = os.WriteFile(filepath, newJson, 0666)
		if err != nil {
			return fmt.Errorf("error creating new database: %w", err)
		}
	} else if err != nil {
		return fmt.Errorf("file error: %w", err)
	}
	return nil
}