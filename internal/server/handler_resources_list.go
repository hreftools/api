package server

import (
	"net/http"

	"github.com/urlspace/api/internal/uow"
)

type resourcesListResponse struct {
	Status string             `json:"status"`
	Data   []responseResource `json:"data"`
}

func handleResourcesList(uowSvc *uow.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := userIDFromContext(r.Context())

		list, err := uowSvc.ListResources(r.Context(), userID)
		if err != nil {
			statusCode, errorMessage := uow.MapErrorToHTTP(err)
			writeJSONError(w, statusCode, errorMessage)
			return
		}

		items := make([]responseResource, len(list))
		for i, item := range list {
			items[i] = newResponseResource(item)
		}

		writeJSONSuccess(w, http.StatusOK, resourcesListResponse{
			Status: "ok",
			Data:   items,
		})
	}
}
