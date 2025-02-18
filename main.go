package main

import (
	"log"
	"net/http"
)



func main() {
	const port  = "8080"
	const filePathRoute = "."
	log.Println("hello world, starting server")

	

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(filePathRoute)))
	mux.Handle("/assets", http.FileServer(http.Dir("./assets/logo.png")))

	server := &http.Server{
		Addr: ":" + port,
		Handler: mux,
	}
	log.Printf("Serving files from %s on port %s\n", filePathRoute, port)
	log.Fatal(server.ListenAndServe())

	
}