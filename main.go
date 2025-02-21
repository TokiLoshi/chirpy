package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func(cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, req)
	})
}

func(cfg *apiConfig) metricsHandler(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf(`<html>
		<body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`, cfg.fileserverHits.Load())))
}

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, req *http.Request) {
	cfg.fileserverHits.Store(0)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

}

func respondWithError(w http.ResponseWriter, code int, message string) {
	type errorResponse struct {
		Error string `json:"error"`
	}
	responseError := errorResponse{
		Error: message,
	}
	data, err := json.Marshal(responseError)
	if err != nil {
		log.Printf("error marshalling JSON %s", err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	dat, err := json.Marshal(payload)
	if err != nil {
		respondWithError(w, 400, "error marshalling JSON")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)	
	w.Write(dat)
}

type cleanedResponse struct {
	CleanedBody string `json:"cleaned_body"`
}

func cleanChirpBody(chirp string) string {
	
	words := strings.Split(chirp, " ")
	cleanedWords := make([]string, 0, len(words))
	forbiddenWords := map[string]bool{
		"kerfuffle": true, 
		"sharbert": true, 
		"fornax": true,
	}
	for _, word := range words {
		loweredWord := strings.ToLower(word)
		if forbiddenWords[loweredWord] {
			cleanedWords = append(cleanedWords, "****")
		} else {
			cleanedWords = append(cleanedWords, word)
		}
		
	}
	cleanChirp := strings.Join(cleanedWords, " ")
	return cleanChirp
}

func validateChirps(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	// Decode the request 
	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	
	// Handle errors decoding the request 
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, "error decoding json data")
		return 
	}

	// Check if chirps length is too long 
	maxLength := 140
	if len(params.Body) > maxLength {
		respondWithError(w, 400, "Chirp is too long")
		return 
	}

	// Check if the chirp is clean 
	cleanchirp := cleanChirpBody(params.Body)
	respondWithJson(w, 200, cleanedResponse{cleanchirp})

}


func main() {
	const port  = "8080"
	const filePathRoute = "."
	apiCfg := &apiConfig{}
	log.Println("hello world, starting server")
	mux := http.NewServeMux()
	
	// Health check endpoints 
	mux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))	
	})

	// Admin endpoints 
	mux.HandleFunc("GET /admin/metrics", apiCfg.metricsHandler)	
	mux.HandleFunc("POST /admin/reset", apiCfg.resetHandler)
	
	// Static file serving
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filePathRoute)))))
	mux.Handle("/assets/", http.FileServer(http.Dir("./assets")))
	
	// API endpoints 
	mux.HandleFunc("POST /api/validate_chirp", validateChirps)

	server := &http.Server{
		Addr: ":" + port,
		Handler: mux,
	}
	log.Printf("Serving files from %s on port %s\n", filePathRoute, port)
	log.Fatal(server.ListenAndServe())

	
}