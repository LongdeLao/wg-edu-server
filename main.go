package main

import (
	"fmt"
	"net/http"
)

func main() {
	// Handle the health check route
	http.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		// Write a 200 OK response with the message "OK"
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Start the server on port 8080
	fmt.Println("Server starting on port 8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
