package middlewares

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/hreftools/api/internal/response"
	"github.com/hreftools/api/internal/store"
	"github.com/hreftools/api/internal/utils"
)

func Admin(s *store.Store) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, ok := utils.UserIDFromContext(r.Context())
			if !ok {
				response.WriteJSONError(w, http.StatusUnauthorized, "unauthorized")
				return
			}

			user, err := s.Users.GetById(r.Context(), userID)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					response.WriteJSONError(w, http.StatusUnauthorized, "unauthorized")
					return
				}
				response.HandleServerError(w, err, "failed to look up user")
				return
			}

			if !user.IsAdmin {
				response.WriteJSONError(w, http.StatusForbidden, "forbidden")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
