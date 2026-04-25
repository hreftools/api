package server

import (
	"encoding/json"
	"net/http"

	"github.com/urlspace/api/internal/uow"
)

type resourceCreateBody struct {
	Title       string   `json:"title"`
	URL         string   `json:"url"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
}

type resourceCreateResponse struct {
	Status string           `json:"status"`
	Data   responseResource `json:"data"`
}

func handleResourcesCreate(uowSvc *uow.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := userIDFromContext(r.Context())

		var body resourceCreateBody
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&body); err != nil {
			handleClientError(w, err, "invalid request body")
			return
		}

		result, err := uowSvc.CreateResource(r.Context(), uow.CreateResourceParams{
			UserID:      userID,
			Title:       body.Title,
			Url:         body.URL,
			Description: body.Description,
			Tags:        body.Tags,
		})
		if err != nil {
			statusCode, errorMessage := uow.MapErrorToHTTP(err)
			writeJSONError(w, statusCode, errorMessage)
			return
		}

		writeJSONSuccess(w, http.StatusCreated, resourceCreateResponse{
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
