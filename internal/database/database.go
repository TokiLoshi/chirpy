package database

import (
	"fmt"
	"sync"
)

type Chirp struct {
	Id int `json:"id"`
	Body string `json:"body"`
}

type DB struct {
	path string 
	mux *sync.RWMutex 
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
}

// NewDB creates a new database connection
// and reates the database file if it doesn't exist 
func NewDB(path string) (*DB, error) {
	fmt.Println("Creating a new db")
	db := &DB{
		path : path,
		mux : &sync.RWMutex{},
	}
	err := db.ensureDB()
	if err != nil {
		return nil, fmt.Errorf("error ensuring the db exists: %w", err)
	}
	return db, nil 
}








