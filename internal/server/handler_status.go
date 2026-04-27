package server

import (
	"net/http"
)

type statusResponse struct {
	Status string `json:"status"`
	Data   string `json:"data"`
}

func handleStatus(w http.ResponseWriter, r *http.Request) {
	writeJSONSuccess(w, http.StatusOK, statusResponse{
		Status: "ok",
		Data:   "service is running",
	})
}
