package main

import (
	"fmt"
	"log"
	"net/http"
)

type apiConfig struct {
	fileserveHits int
}

func (c *apiConfig) reset() {
	c.fileserveHits = 0
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
	mux.HandleFunc("/healthz", readinessHandler)
	mux.HandleFunc("/metrics", apiCfg.hitHandler)
	mux.HandleFunc("/reset", apiCfg.resetHandler)
	

	srv := &http.Server{
		Addr: ":" + port,
		Handler: mux,
	}
	
	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
	
}

func readinessHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))

}

func (c *apiConfig) hitHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hits: %d", c.fileserveHits)
}

func (c *apiConfig) resetHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("resetting fileserveHits")
	c.reset()
	w.Header().Set("Content-Type", "text/lain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Status reset to %d", c.fileserveHits)
}