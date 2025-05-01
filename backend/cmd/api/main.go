// main is the entry point for the API.
package main

import (
	"context"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/edsonjaramillo/hyde/backend/internal/aes"
	"github.com/edsonjaramillo/hyde/backend/internal/env"
)

func getFile(r *http.Request) ([]byte, error) {
	file, _, err := r.FormFile("file")
	if err != nil {
		return nil, errors.New("error retrieving file")
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("Error closing file: %v", err)
		}
	}()

	if file == nil {
		return nil, errors.New("file is empty")
	}

	inputFileStream, err := io.ReadAll(file)
	if err != nil {
		return nil, errors.New("error reading file")
	}

	return inputFileStream, nil
}

func writeResponse(w http.ResponseWriter, data []byte) {
	w.Header().Set("Content-Type", "application/octet-stream")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write(data)
	if err != nil {
		log.Printf("Error writing response: %v", err)
		http.Error(w, "Error writing response", http.StatusInternalServerError)
		return
	}
}

func encryptHandler(w http.ResponseWriter, r *http.Request) {
	filestream, err := getFile(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	encryptedData, err := aes.EncryptData(filestream, r.FormValue("password"))
	if err != nil {
		http.Error(w, "Error encrypting.", http.StatusInternalServerError)
		return
	}

	writeResponse(w, encryptedData)
}

func decryptHandler(w http.ResponseWriter, r *http.Request) {
	filestream, err := getFile(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	decryptedData, err := aes.DecryptData(filestream, r.FormValue("password"))
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusInternalServerError)
		return
	}

	writeResponse(w, decryptedData)
}

func corsMiddleware(allowedOrigin string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Set CORS headers for all responses
			w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
			w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type") // Add others if needed (e.g., Authorization)
			w.Header().Set("Access-Control-Max-Age", "86400")              // 1 day

			// Handle preflight requests
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}

			// Call the next handler in the chain
			next.ServeHTTP(w, r)
		})
	}
}

func main() {
	// --- Environment Setup ---
	env := env.Values

	// --- Router Setup ---
	mux := http.NewServeMux()

	// Apply middleware chain (CORS)
	corsHandler := corsMiddleware(env.AllowedOrgins)
	finalEncryptHandler := corsHandler(http.HandlerFunc(encryptHandler))
	finalDecryptHandler := corsHandler(http.HandlerFunc(decryptHandler))

	// Register handlers
	mux.Handle("POST /encrypt", finalEncryptHandler)
	mux.Handle("POST /decrypt", finalDecryptHandler)

	// --- Server Setup ---
	server := &http.Server{
		Addr:    ":" + env.APIPort,
		Handler: mux,
	}

	// --- Graceful Shutdown ---
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("Server listening on port :%s", env.APIPort)
		log.Printf("Allowed Origin: %s", env.AllowedOrgins)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Could not listen on %s: %v\n", server.Addr, err)
		}
	}()

	// Wait for shutdown signal
	<-stopChan
	log.Println("Shutting down server...")

	// Give ongoing requests a deadline to finish
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}

	log.Println("Server exited properly")
}
