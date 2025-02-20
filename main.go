package main

import (
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
		hits := fmt.Sprintf("Hits: %d", cfg.fileserverHits.Load())
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(hits))
}

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, req *http.Request) {
	cfg.fileserverHits.Store(0)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}


func main() {
	const port  = "8080"
	const filePathRoute = "."
	apiCfg := &apiConfig{}
	log.Println("hello world, starting server")

	

	mux := http.NewServeMux()
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filePathRoute)))))
	mux.Handle("/assets/", http.FileServer(http.Dir("./assets")))
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))	
	})
	mux.HandleFunc("/metrics", apiCfg.metricsHandler)	
	mux.HandleFunc("/reset", apiCfg.resetHandler)

	server := &http.Server{
		Addr: ":" + port,
		Handler: mux,
	}
	log.Printf("Serving files from %s on port %s\n", filePathRoute, port)
	log.Fatal(server.ListenAndServe())

	
}