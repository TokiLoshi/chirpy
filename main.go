package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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

func validateChirps(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type errorResponse struct {
		Error string `json:"error"`
	}

	// Decode the request 
	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	
	// Handle errors decoding the request 
	err := decoder.Decode(&params)
	if err != nil {
		decodeError := errorResponse{
			Error : "error decoding json data",
		}
		data, error := json.Marshal(decodeError)
		if error != nil {
			log.Printf("error marshalling JSON %s", err)
		}
		log.Printf("Error decoding parameters %s", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		w.Write(data)
		return 
	}

	// Check if chirps length is too long 
	if len(params.Body) > 140 {
		errorRespBody := errorResponse{
			Error: "Chirp is too long",
		}
		data, err := json.Marshal(errorRespBody)
		if err != nil {
			log.Printf("error marshalling JSON: %s", err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		w.Write(data)
		return 
	}
	type correctResponse struct {
		Valid bool `json:"valid"`
	}
	respBody := correctResponse{
		Valid : true,
	}
	dat, err := json.Marshal(respBody)
	if err != nil {
		log.Printf("error marshalling JSON: %s", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)	
	w.Write(dat)

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