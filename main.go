package main

import (
	"log"
	"net/http"
)

type apiConfig struct {
	fileserveHits int
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

	// Create new multiplexer 
	mux := http.NewServeMux()
	// File Server to serve root path 
	fileServer := http.FileServer(http.Dir(filepathRoot))
	// Handle root and strip out the last '/'
	handler :=  http.StripPrefix("/app/", fileServer)
	// Handle the root with the middleware 
	apiCfg := apiConfig{}
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(handler))
	// Handle Readiness function
	mux.HandleFunc("GET /api/healthz", readinessHandler)
	mux.HandleFunc("GET /admin/metrics", apiCfg.hitHandler)
	mux.HandleFunc("GET /api/reset", apiCfg.resetHandler)
	mux.HandleFunc("POST /api/validate_chirp", apiCfg.validateChirp)
	

	srv := &http.Server{
		Addr: ":" + port,
		Handler: mux,
	}
	
	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
	
}





