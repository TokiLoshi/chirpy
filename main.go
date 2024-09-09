package main

import (
	"fmt"
	"net/http"
)

type Server struct{}

func (Server) ServeHTTP(http.ResponseWriter, *http.Request) {}

func main() {

	// 1. Create a new http.ServeMux
	mux := http.NewServeMux()
	mux.Handle("/localhost", Server{})
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		// The "/" pattern matches everything, so we need to check 
		// that we're at the root here 
		if req.URL.Path != "/ok" {
			http.NotFound(w, req)
			return
		}
		fmt.Fprintf(w, "Home Page")
	})

	http.ListenAndServe(":8080", mux)

}