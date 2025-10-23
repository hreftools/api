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

func main() {
	myHandler := http.NewServeMux()

	myHandler.HandleFunc("GET /status", func(w http.ResponseWriter, r *http.Request) {
		response := ResponseSuccess{
			Status: "success",
			Data:   "Service is running",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	s := &http.Server{
		Addr:           port,
		Handler:        myHandler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Printf("Starting server on %s", port)

	log.Fatal(s.ListenAndServe())
}
