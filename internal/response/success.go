package response

import (
	"encoding/json"
	"log"
	"net/http"
)

func WriteJSONSuccess(w http.ResponseWriter, statusCode int, res any) {
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Printf("Error encoding error response: %v", err)
	}
}
