package server

import (
	"net/http"

	"github.com/urlspace/api/internal/tag"
)

type tagsListResponse struct {
	Status string        `json:"status"`
	Data   []responseTag `json:"data"`
}

func handleTagsList(svc *tag.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := userIDFromContext(r.Context())

		list, err := svc.List(r.Context(), userID)
		if err != nil {
			statusCode, errorMessage := tag.MapErrorToHTTP(err)
			writeJSONError(w, statusCode, errorMessage)
			return
		}

		items := make([]responseTag, len(list))
		for i, item := range list {
			items[i] = responseTag{
				ID:        item.ID,
				Name:      item.Name,
				CreatedAt: item.CreatedAt,
				UpdatedAt: item.UpdatedAt,
			}
		}

		writeJSONSuccess(w, http.StatusOK, tagsListResponse{
			Status: "ok",
			Data:   items,
		})
	}
}
