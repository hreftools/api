package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

const (
	port = ":8080"
)

type ResponseSuccess struct {
	Status string `json:"status"`
	Data   string `json:"data"`
}

type ResponseError struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

func commonHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Security headers
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Content-Security-Policy", "default-src 'self'")

		// CORS
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Others
		w.Header().Set("Content-Type", "application/json")

		next.ServeHTTP(w, r)
	})
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /status", func(w http.ResponseWriter, r *http.Request) {
		response := ResponseSuccess{
			Status: "success",
			Data:   "Service is running",
		}

		json.NewEncoder(w).Encode(response)
	})

	handler := commonHeaders(mux)

	s := &http.Server{
		Addr:              port,
		Handler:           handler,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	log.Printf("Starting server on %s", port)

	log.Fatal(s.ListenAndServe())
}
