package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/TokiLoshi/chirpy/internal/database"
	"github.com/joho/godotenv"
)

type apiConfig struct {
	fileserveHits int
	DB *database.DB
	jwt []byte
}



func (c *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.fileserveHits++
		log.Println("File server hit")
		next.ServeHTTP(w, r)
	})
		
}


func main() {
	const port = "8080"
	const filepathRoot = "."
	db, err := database.NewDB("./databse.json")
	if err != nil {
		log.Fatalf("Failed to initialize the database: %v", err)
	}

	envError := godotenv.Load()
	if envError != nil {
		log.Fatal("Error loading env file")
	}
	secret := os.Getenv("JWT_SECRET")
	fmt.Printf("Secret: %v\n", secret)

	// Create new multiplexer 
	mux := http.NewServeMux()
	// File Server to serve root path 
	fileServer := http.FileServer(http.Dir(filepathRoot))
	// Handle root and strip out the last '/'
	handler :=  http.StripPrefix("/app/", fileServer)
	// Handle the root with the middleware 
	apiCfg := apiConfig{
		DB: db,
		jwt: []byte(secret),
	}
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(handler))
	// Handle Readiness function
	mux.HandleFunc("/api/healthz", readinessHandler)
	mux.HandleFunc("/admin/metrics", apiCfg.hitHandler)
	mux.HandleFunc("/api/reset", apiCfg.resetHandler)
	// mux.HandleFunc("POST /api/validate_chirp", apiCfg.validateChirp)
	mux.HandleFunc("/api/chirps", apiCfg.createChripHandler)
	mux.HandleFunc("/api/chirps/{chirpId}", apiCfg.singleChirpHandler)
	
	mux.HandleFunc("POST /api/login", apiCfg.userLoginHandler)
	mux.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			apiCfg.createUserHandler(w, r) 
		} else if r.Method == http.MethodPut {
			apiCfg.updateUserHandler(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	srv := &http.Server{
		Addr: ":" + port,
		Handler: mux,
	}
	
	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
	
}





