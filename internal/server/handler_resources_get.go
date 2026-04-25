package server

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/urlspace/api/internal/uow"
)

type resourcesGetResponse struct {
	Status string           `json:"status"`
	Data   responseResource `json:"data"`
}

func handleResourcesGet(uowSvc *uow.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := userIDFromContext(r.Context())

		id := r.PathValue("id")
		idUuid, err := uuid.Parse(id)
		if err != nil {
			handleClientError(w, err, "invalid id parameter")
			return
		}

		result, err := uowSvc.GetResource(r.Context(), idUuid, userID)
		if err != nil {
			statusCode, errorMessage := uow.MapErrorToHTTP(err)
			writeJSONError(w, statusCode, errorMessage)
			return
		}

		writeJSONSuccess(w, http.StatusOK, resourcesGetResponse{
			Status: "ok",
			Data: responseResource{
				ID:          result.ID,
				Title:       result.Title,
				Description: result.Description,
				URL:         result.URL,
				Tags:        result.Tags,
				CreatedAt:   result.CreatedAt,
				UpdatedAt:   result.UpdatedAt,
			},
		})
	}
}
