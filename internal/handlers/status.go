package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/zapi-sh/api/internal/models"
)

func Status(w http.ResponseWriter, r *http.Request) {

	response := models.ResponseSuccess{
		Status: "success",
		Data:   "Service is running",
	}

	json.NewEncoder(w).Encode(response)
}
