package main

import (
	"log"
	"net/http"
)

// func (serverHandler) ServeHTTP(http.ResponseWriter, *http.Request) {

// }

func main() {
	const port = "8080" 
	mux := http.NewServeMux()

	server := &http.Server{
		Addr: ":" + port,
		Handler: mux,
	}
	
	fileServer := http.FileServer(http.Dir("."))
	mux.Handle("/", fileServer)

	log.Printf("Serving on port: %s\n", port)
	http.ListenAndServe(server.Addr, server.Handler)
	
	

}