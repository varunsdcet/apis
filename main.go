package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Struct for JSON response
type Response struct {
	Message string `json:"message,omitempty"`
	Time    string `json:"time,omitempty"`
}

// Root handler for "/"
func rootHandler(w http.ResponseWriter, r *http.Request) {
	response := Response{Message: "Welcome to the Go server!"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Hello handler for "/hello"
func helloHandler(w http.ResponseWriter, r *http.Request) {
	response := Response{Message: "Hello, Golang!"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Time handler for "/time"
func timeHandler(w http.ResponseWriter, r *http.Request) {
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	response := Response{Time: currentTime}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/time", timeHandler)

	fmt.Println("Server starting at port 8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
