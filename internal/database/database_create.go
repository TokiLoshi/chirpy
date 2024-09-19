package database

import (
	"fmt"
	"strings"
)

func (db *DB) CreateUser(body string) (User, error) {
	fmt.Println("Creating a new user")
	db.mux.Lock()
	defer db.mux.Unlock()

	newUser := User{}

	email := body 
	dbStructure, err := db.LoadDB()
	if err != nil {
		return User{}, fmt.Errorf("error loading db: %w", err)
	}
	fmt.Printf("Finished loading db: %v\n", dbStructure)

	newId := highestUserId(dbStructure) + 1 
	fmt.Printf("Assigning new user: %v\n", newId)

	newUser = User {
		Id: newId, 
		Email: email, 
	}
	fmt.Printf("New user: %v\n", newUser)

	dbStructure.Users[newUser.Id] = newUser 

	fmt.Printf("adding user to db struct: %v\n", dbStructure)

	err = db.WriteDB(dbStructure)
	if err != nil {
		return User{}, fmt.Errorf("error writing user to db: %v", err)
	}
	fmt.Println("finished writing user to db ")

	return newUser, nil 
}

// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateChirp(body string) (Chirp, error) {
	fmt.Println("Creating a new chirp")
	db.mux.Lock()
	defer db.mux.Unlock()

	newChirp := Chirp{}


	// Validate and clean the chirp 
	validChirp, err := validateChirp(body)
	if err != nil {
		return newChirp, fmt.Errorf("error validating chirp: %w", err)
	}
	fmt.Printf("Returned from validating chirp: %v\n", validChirp)

	words := strings.Split(validChirp, " ")
	var cleanedWords []string 
	for _, word := range words {
		cleanedWords = append(cleanedWords, cleanProfane(word))
	}
	cleanedChirp := strings.Join(cleanedWords, " ")
	fmt.Printf("Finished cleaning chirp: %v\n", cleanedChirp)

	dbStructure, err := db.LoadDB()
	if err != nil {
		return Chirp{}, fmt.Errorf("error loading db: %w", err)
	}
	fmt.Printf("Finished loading db: %v\n", dbStructure)
	
	newId := highestId(dbStructure) + 1
	fmt.Printf("Assigning new id: %v\n", newId)

	newChirp = Chirp{
		Id: newId,
		Body: cleanedChirp,
	}

	dbStructure.Chirps[newChirp.Id] = newChirp

	fmt.Printf("adding chirp to db struct: %v\n", dbStructure)

	err = db.WriteDB(dbStructure) 
	if err != nil {
		fmt.Printf("This is where things are going horribly wrong in the writing, it never returns")
		return Chirp{}, fmt.Errorf("error writing db: %w", err)
	}

	fmt.Println("Finished writing to DB. And we're all done")
	

	// Write the chirp 
	return newChirp, nil
}

func highestUserId(users DBStructure) int {
	highestId := 0 
	for _, user := range users.Users {
		if user.Id > highestId {
			highestId = user.Id
		}
	}
	return highestId
}

func highestId(chirps DBStructure) int {
	highestId := 0
	for _, chirp := range chirps.Chirps {
		if chirp.Id > highestId {
			highestId = chirp.Id
		}
	}
	fmt.Printf("Returning highest id: %v", highestId)
	return highestId
}

func validateChirp(chirp string) (string, error) {
	// check the length 
	if len(chirp) > 140 {
		return "", fmt.Errorf("chirp is too long")
	}
	fmt.Printf("Returning valid chirp: %v", chirp)
	return chirp, nil
}


func cleanProfane(chirp string) string {
	profane := strings.ToLower(chirp)
	cleanChirp := profane
	switch profane {
	case "kerfuffle":
		cleanChirp = "****"
	case "sharbert":
		cleanChirp = "****"
	case "fornax":
		cleanChirp = "****"
	default:
		cleanChirp = chirp
	}
	fmt.Printf("Returning less profane chirps: %v", cleanChirp)
	return cleanChirp
}