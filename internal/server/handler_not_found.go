package server

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type notFoundResponse struct {
	Status string `json:"status"`
	Data   string `json:"data"`
}

func handleNotFound(w http.ResponseWriter, r *http.Request) {
	response := &notFoundResponse{
		Status: "error",
		Data:   "endpoint not found",
	}

	w.WriteHeader(http.StatusNotFound)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		slog.Error("failed to encode error response", "error", err)
	}
}
