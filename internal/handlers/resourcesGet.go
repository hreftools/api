package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/zapi-sh/api/internal/models"
)

func ResourcesGet(w http.ResponseWriter, r *http.Request) {

	response := models.ResponseSuccess{
		Status: "ok",
		Data:   "data goes here",
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
