package server

import (
	"net/http"

	"github.com/urlspace/api/internal/user"
)

type meGetResponse struct {
	Status string       `json:"status"`
	Data   responseUser `json:"data"`
}

func handleMeGet(svc *user.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := userIDFromContext(r.Context())
		if !ok {
			writeJSONError(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		u, err := svc.GetById(r.Context(), userID)
		if err != nil {
			statusCode, errorMessage := user.MapErrorToHTTP(err)
			writeJSONError(w, statusCode, errorMessage)
			return
		}

		writeJSONSuccess(w, http.StatusOK, meGetResponse{
			Status: "ok",
			Data:   newResponseUser(u),
		})
	}
}
