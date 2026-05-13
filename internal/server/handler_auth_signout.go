package server

import (
	"net/http"

	"github.com/urlspace/api/internal/user"
)

type authSignoutResponse struct {
	Status string `json:"status"`
	Data   string `json:"data"`
}

func handleAuthSignout(svc *user.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := resolveSession(r)

		if err := svc.Signout(r.Context(), session); err != nil {
			statusCode, errorMessage := user.MapErrorToHTTP(r.Context(), err)
			writeJSONError(w, statusCode, errorMessage)
			return
		}

		clearSessionCookie(w)

		writeJSONSuccess(w, http.StatusOK, authSignoutResponse{
			Status: "ok",
			Data:   "ok",
		})
	}
}
