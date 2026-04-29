package server

import "net/http"

func handleNotFound(w http.ResponseWriter, r *http.Request) {
	writeJSONError(w, http.StatusNotFound, "endpoint not found")
}
