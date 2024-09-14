package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
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
	mux.HandleFunc("GET /api/healthz", readinessHandler)
	mux.HandleFunc("GET /admin/metrics", apiCfg.hitHandler)
	mux.HandleFunc("GET /api/reset", apiCfg.resetHandler)
	

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
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	
	hits := c.fileserveHits
	data, err := os.ReadFile("./admin.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("We got an error reading the admin html: %v", err)
		fmt.Fprint(w, "Internal Server Error")
		return
	}
	w.WriteHeader(http.StatusOK)
	stringifiedData := string(data)
	htmlContent := fmt.Sprintf(stringifiedData, hits)
	fmt.Fprint(w, htmlContent)

}

func (c *apiConfig) resetHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("resetting fileserveHits")
	c.reset()
	w.Header().Set("Content-Type", "text/lain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Status reset to %d", c.fileserveHits)
}