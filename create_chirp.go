package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

func (c *apiConfig) createChirp(w http.ResponseWriter, r *http.Request) {
	
	// define the struct parameters 
	type parameters struct {
		Body string `json:"body"`
	}

	type returnErr struct {
		Error string `json:"error"` 
	}

	type okStruct struct {
		Valid string `json:"cleaned_body"`
	}

	// define a new decoder and params 
	
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		fmt.Printf("error: Something went wrong")
		return
	}
	
	var dat []byte

	// Handle chirps that are too long 
	if len(params.Body) > 140 {
		w.WriteHeader(400)
		dat, err = json.Marshal(returnErr{Error: "Chirp is too long"})
	} else {
		words := strings.Split(params.Body, " ")
		var cleanedWords []string
		for _, word := range words {
			cleanedWords = append(cleanedWords, cleanProfane(word))
		}
		newChirp := strings.Join(cleanedWords, " ")

		// Open the json file 
		file, err := os.Open("./database.json")
		if err != nil {
			fmt.Printf("error opening json file: %v\n", err)
		}
		defer file.Close()

		// Read in data 
		data := make([]byte, 100)
		count, err := file.Read(data)
		if err != nil {
			fmt.Printf("error reading json file: %v\n", err)
		}
		fmt.Printf("read: %d bytes: %q\n", count, data[:count])
		parameters.Id = count + 1

		w.WriteHeader(200)
		dat, err = json.Marshal(okStruct{Valid:newChirp})
	}

	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}

	//Chirp is ok we can return it 
	
	// Assign a unique id (we're just incrementing chirps by 1)
	// Save it to disk 
	// if all goes well respond with a 201 status code 
	// return full chirp resource 
	w.Header().Set("Content-Type", "application/json")
	w.Write(dat)
}

// func cleanProfane(chirp string) string {
// 	profane := strings.ToLower(chirp)
// 	cleanChirp := profane
// 	switch profane {
// 	case "kerfuffle":
// 		cleanChirp = "****"
// 	case "sharbert":
// 		cleanChirp = "****"
// 	case "fornax":
// 		cleanChirp = "****"
// 	default:
// 		cleanChirp = chirp
// 	}
// 	return cleanChirp
// }