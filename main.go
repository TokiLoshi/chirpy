package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/TokiLoshi/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)


type apiConfig struct {
	fileserverHits atomic.Int32
	dbQueries *database.Queries
}


func main() {

	// Load the env variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Open Database connection
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to open DB: %v", err)
	}
	dbQueries := database.New(db)

	const port  = "8080"
	const filePathRoute = "."
	apiCfg := &apiConfig{
		dbQueries: dbQueries,
	}
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