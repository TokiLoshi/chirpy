package main

import (
	"fmt"
	"net/http"
	"os"
)

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