package server

import (
	"net/http"

	"github.com/urlspace/api/internal/collection"
)

type collectionsListResponse struct {
	Status string               `json:"status"`
	Data   []responseCollection `json:"data"`
}

func handleCollectionsList(collectionSvc *collection.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := userIDFromContext(r.Context())

		list, err := collectionSvc.List(r.Context(), userID)
		if err != nil {
			statusCode, errorMessage := collection.MapErrorToHTTP(err)
			writeJSONError(w, statusCode, errorMessage)
			return
		}

		items := make([]responseCollection, len(list))
		for i, item := range list {
			items[i] = newResponseCollection(item)
		}

		writeJSONSuccess(w, http.StatusOK, collectionsListResponse{
			Status: "ok",
			Data:   items,
		})
	}
}
