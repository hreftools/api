package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type ErrorResponse struct {
	Status string `json:"status"`
	Data   string `json:"data"`
}

func HandleError(w http.ResponseWriter, err error) bool {
	if err == nil {
		return false
	}

	var data string
	var statusCode int

	if errors.Is(err, sql.ErrNoRows) {
		statusCode = http.StatusNotFound
		data = "resource not found"
	} else if errors.Is(err, context.DeadlineExceeded) {
		statusCode = http.StatusInternalServerError
		data = "request timeout"
	} else if errors.Is(err, context.Canceled) {
		statusCode = 499
		data = "request cancelled"
	} else {
		statusCode = http.StatusInternalServerError
		data = "internal server error"
		log.Printf("Internal server error: %v", err)
	}

	w.WriteHeader(statusCode)
	response := &ErrorResponse{
		Status: "error",
		Data:   data,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding error response: %v", err)
	}

	return true
}
