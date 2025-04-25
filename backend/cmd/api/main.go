// main is the entry point for the API.
package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	encryption "github.com/edsonjaramillo/hyde/backend/internal/encryption"
)

func encryptHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if file == nil {
		http.Error(w, "File is required", http.StatusBadRequest)
		return
	}

	inputFileStream, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	encryptedData, err := encryption.EncryptFile(inputFileStream, r.FormValue("password"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(encryptedData)
}

func decryptHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if file == nil {
		http.Error(w, "File is required", http.StatusBadRequest)
		return
	}

	content, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	decryptedData, err := encryption.DecryptFile(content, r.FormValue("password"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(decryptedData)
}

// Handler for the OPTIONS preflight request
func preflightHandler(w http.ResponseWriter, r *http.Request) {
	// These headers tell the browser what is allowed for the actual request
	w.Header().Set("Access-Control-Allow-Origin", "*")              // Allow any origin
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS") // Allow POST and OPTIONS
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")  // Allow the Content-Type header
	w.Header().Set("Access-Control-Max-Age", "86400")               // Optional: how long the preflight result can be cached

	// Respond with 200 OK to indicate permission
	w.WriteHeader(http.StatusOK)
}

func main() {
	router := http.NewServeMux()
	router.HandleFunc("POST /encrypt", encryptHandler)
	router.HandleFunc("OPTIONS /encrypt", preflightHandler)

	router.HandleFunc("POST /decrypt", decryptHandler)
	router.HandleFunc("OPTIONS /decrypt", preflightHandler)

	server := http.Server{Addr: ":8080", Handler: router}

	fmt.Println("Server listening on port :8080")
	err := server.ListenAndServe()
	if err != nil {
		log.Println("Error starting server:", err)
	}
}
