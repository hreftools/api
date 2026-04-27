package server

import (
	"encoding/json"
	"net/http"

	"github.com/urlspace/api/internal/collection"
)

type collectionCreateBody struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Public      bool   `json:"public"`
}

type collectionCreateResponse struct {
	Status string             `json:"status"`
	Data   responseCollection `json:"data"`
}

func handleCollectionsCreate(collectionSvc *collection.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := userIDFromContext(r.Context())

		var body collectionCreateBody
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&body); err != nil {
			handleClientError(w, err, "invalid request body")
			return
		}

		result, err := collectionSvc.Create(r.Context(), collection.CreateParams{
			UserID:      userID,
			Title:       body.Title,
			Description: body.Description,
			Public:      body.Public,
		})
		if err != nil {
			statusCode, errorMessage := collection.MapErrorToHTTP(err)
			writeJSONError(w, statusCode, errorMessage)
			return
		}

		writeJSONSuccess(w, http.StatusCreated, collectionCreateResponse{
			Status: "ok",
			Data:   newResponseCollection(result),
		})
	}
}
