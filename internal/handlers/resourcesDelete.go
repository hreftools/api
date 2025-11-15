package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/zapi-sh/api/internal/db"
	"github.com/zapi-sh/api/internal/store"
)

ResourceDeleteResponse struct {
	Status string        `json:"status"`
	Data   db.Resource `json:"data"`
}

func ResourcesDelete(store *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		list, err := store.Resources(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		s := make([]db.Resource, 0, len(list))
		for _, item := range list {
			s = append(s, item)
		}

		response := &ResourcesGetResponse{
			Status: "ok",
			Data:   s,
		}

		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
