package middlewares

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hreftools/api/internal/config"
	"github.com/hreftools/api/internal/response"
	"github.com/hreftools/api/internal/store"
)

func Auth(s *store.Store) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenID, ok := resolveTokenID(r)
			if !ok {
				response.WriteJSONError(w, http.StatusUnauthorized, "unauthorized")
				return
			}

			token, err := s.Tokens.GetByID(r.Context(), tokenID)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					response.WriteJSONError(w, http.StatusUnauthorized, "unauthorized")
					return
				}
				response.HandleServerError(w, err, "failed to look up token")
				return
			}

			if time.Now().After(token.ExpiresAt) {
				response.WriteJSONError(w, http.StatusUnauthorized, "unauthorized")
				return
			}

			// Sliding expiry: renew session tokens that are approaching expiry.
			if token.Type == config.TokenTypeSession && time.Until(token.ExpiresAt) < config.SessionRenewalThreshold {
				go func() {
					// Fire-and-forget renewal. Errors are intentionally swallowed:
					// a failed renewal is non-fatal — the token remains valid for
					// the remainder of its current expiry window and renewal will
					// be retried on the next request.
					_, _ = s.Tokens.UpdateExpiresAt(context.Background(), store.TokenUpdateExpiresAtParams{
						ID:        token.ID,
						ExpiresAt: time.Now().Add(config.SessionExpiryDuration),
					})
				}()
			}

			ctx := context.WithValue(r.Context(), config.UserIDContextKey, token.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// resolveTokenID extracts the token UUID from the Authorization header (Bearer scheme)
// or falls back to the session_id cookie.
func resolveTokenID(r *http.Request) (uuid.UUID, bool) {
	if authHeader := r.Header.Get("Authorization"); authHeader != "" {
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
			if id, err := uuid.Parse(strings.TrimSpace(parts[1])); err == nil {
				return id, true
			}
		}
	}

	if cookie, err := r.Cookie(config.SessionCookieName); err == nil {
		if id, err := uuid.Parse(cookie.Value); err == nil {
			return id, true
		}
	}

	return uuid.UUID{}, false
}
