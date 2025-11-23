package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/zapi-sh/api/internal/models"
)

func Status(w http.ResponseWriter, r *http.Request) {

	response := models.ResponseSuccess{
		Status: "success",
		Data:   "service is running",
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
