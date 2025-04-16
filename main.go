package main

import (
	"fmt"
	"net/http"
	"time"
)

// Hello handler for the "/hello" route
func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, Golang!")
}

// Time handler for the "/time" route
func timeHandler(w http.ResponseWriter, r *http.Request) {
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	fmt.Fprintf(w, "Current time: %s", currentTime)
}

func main() {
	// Handle the "/hello" route
	http.HandleFunc("/hello", helloHandler)

	// Handle the "/time" route
	http.HandleFunc("/time", timeHandler)

	// Start the server on port 8080
	fmt.Println("Server starting at port 8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
