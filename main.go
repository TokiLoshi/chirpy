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
	jwtSecret string
	polkaKey string
}


func main() {

	// Load the env variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Open Database connection
	dbURL := os.Getenv("DB_URL")
	log.Printf("DB URL: %s", dbURL)
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to open DB: %v", err)
	}

	// Ping to check the database connection 
	err = db.Ping() 
	if err != nil {
		log.Fatalf("Failed to connect to db DB: %v", err)
	}
	log.Println("Successfully connected to db")

	// Load jwt secret ont apiConfig 
	jwtSecretString := os.Getenv("JWT_SECRET")
	if len(jwtSecretString) == 0 {
		log.Printf("Failed to extract jwt: %v", jwtSecretString)
	} 

	// Load polka api key onto apiConfig 
	polkaKeyString := os.Getenv("POLKA_KEY")
	if len(polkaKeyString) == 0 {
		log.Printf("Failed to extract polka key: %v", polkaKeyString)
	}

	dbQueries := database.New(db)

	const port  = "8080"
	const filePathRoute = "."
	apiCfg := &apiConfig{
		dbQueries: dbQueries,
		jwtSecret: jwtSecretString,
		polkaKey: polkaKeyString,
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
	mux.HandleFunc("POST /admin/reset-metrics", apiCfg.resetHandler)
	mux.HandleFunc("POST /admin/reset", apiCfg.adminReset)
	
	// Static file serving
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filePathRoute)))))
	mux.Handle("/assets/", http.FileServer(http.Dir("./assets")))
	
	// User API endpoints 
	mux.HandleFunc("POST /api/users", apiCfg.validateUser)
	mux.HandleFunc("POST /api/login", apiCfg.handleLogin)
	mux.HandleFunc("POST /api/refresh", apiCfg.hanldeRefresh)
	mux.HandleFunc("POST /api/revoke", apiCfg.handleRevoke)
	mux.HandleFunc("PUT /api/users", apiCfg.handleAuthentication)

	// Chirp API endpoints 
	mux.HandleFunc("POST /api/chirps", apiCfg.createChirp)
	mux.HandleFunc("GET /api/chirps", apiCfg.getAllChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.getSingleChirp)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.handleDelete)
	
	// web hooks 
	mux.HandleFunc("POST /api/polka/webhooks", apiCfg.handleUpgrade)

	server := &http.Server{
		Addr: ":" + port,
		Handler: mux,
	}
	log.Printf("Serving files from %s on port %s\n", filePathRoute, port)
	log.Fatal(server.ListenAndServe())

	
}