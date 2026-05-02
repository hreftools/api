package server

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/urlspace/api/internal/collection"
	"github.com/urlspace/api/internal/config"
	"github.com/urlspace/api/internal/uow"
	"github.com/urlspace/api/internal/user"
)

type responseLinkCollection struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type responseLinkTag struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type responseLink struct {
	ID          uuid.UUID               `json:"id"`
	Title       string                  `json:"title"`
	Description string                  `json:"description"`
	URL         string                  `json:"url"`
	Tags        []responseLinkTag       `json:"tags"`
	Collection  *responseLinkCollection `json:"collection"`
	CreatedAt   time.Time               `json:"createdAt"`
	UpdatedAt   time.Time               `json:"updatedAt"`
}

type responseCollection struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Public      bool      `json:"public"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func newResponseLink(r uow.EnrichedLink) responseLink {
	var col *responseLinkCollection
	if r.Collection != nil {
		col = &responseLinkCollection{ID: r.Collection.ID, Name: r.Collection.Name}
	}
	tags := make([]responseLinkTag, len(r.Tags))
	for i, t := range r.Tags {
		tags[i] = responseLinkTag{ID: t.ID, Name: t.Name}
	}
	return responseLink{
		ID:          r.ID,
		Title:       r.Title,
		Description: r.Description,
		URL:         r.URL,
		Tags:        tags,
		Collection:  col,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}
}

func newResponseCollection(c collection.Collection) responseCollection {
	return responseCollection{
		ID:          c.ID,
		Name:        c.Name,
		Description: c.Description,
		Public:      c.Public,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
	}
}

type responseTag struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type responseUser struct {
	ID          uuid.UUID `json:"id"`
	Email       string    `json:"email"`
	Username    string    `json:"username"`
	DisplayName string    `json:"displayName"`
	IsAdmin     bool      `json:"isAdmin"`
	IsPro       bool      `json:"isPro"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func newResponseUser(u user.User) responseUser {
	return responseUser{
		ID:          u.ID,
		Email:       u.Email,
		Username:    u.Username,
		DisplayName: u.DisplayName,
		IsAdmin:     u.IsAdmin,
		IsPro:       u.IsPro,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
	}
}

type responseUserAdmin struct {
	ID                              uuid.UUID     `json:"id"`
	Email                           string        `json:"email"`
	EmailVerified                   bool          `json:"emailVerified"`
	EmailVerificationToken          uuid.NullUUID `json:"emailVerificationToken"`
	EmailVerificationTokenExpiresAt *time.Time    `json:"emailVerificationTokenExpiresAt"`
	PasswordResetToken              uuid.NullUUID `json:"passwordResetToken"`
	PasswordResetTokenExpiresAt     *time.Time    `json:"passwordResetTokenExpiresAt"`
	Username                        string        `json:"username"`
	DisplayName                     string        `json:"displayName"`
	IsAdmin                         bool          `json:"isAdmin"`
	IsPro                           bool          `json:"isPro"`
	CreatedAt                       time.Time     `json:"createdAt"`
	UpdatedAt                       time.Time     `json:"updatedAt"`
}

func newResponseUserAdmin(u user.User) responseUserAdmin {
	return responseUserAdmin{
		ID:                              u.ID,
		Email:                           u.Email,
		EmailVerified:                   u.EmailVerified,
		EmailVerificationToken:          u.EmailVerificationToken,
		EmailVerificationTokenExpiresAt: u.EmailVerificationTokenExpiresAt,
		PasswordResetToken:              u.PasswordResetToken,
		PasswordResetTokenExpiresAt:     u.PasswordResetTokenExpiresAt,
		Username:                        u.Username,
		DisplayName:                     u.DisplayName,
		IsAdmin:                         u.IsAdmin,
		IsPro:                           u.IsPro,
		CreatedAt:                       u.CreatedAt,
		UpdatedAt:                       u.UpdatedAt,
	}
}

type responseToken struct {
	ID          uuid.UUID `json:"id"`
	Description string    `json:"description"`
	LastUsedAt  time.Time `json:"lastUsedAt"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func newResponseToken(t user.Token) responseToken {
	return responseToken{
		ID:          t.ID,
		Description: t.Description,
		LastUsedAt:  t.LastUsedAt,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
}

// Request helpers

func resolveSessionID(r *http.Request) (uuid.UUID, bool) {
	if cookie, err := r.Cookie(config.SessionCookieName); err == nil {
		if id, err := uuid.Parse(cookie.Value); err == nil {
			return id, true
		}
	}
	return uuid.UUID{}, false
}

func resolveBearerToken(r *http.Request) (string, bool) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", false
	}
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return "", false
	}
	token := strings.TrimSpace(parts[1])
	if token == "" {
		return "", false
	}
	return token, true
}

func userIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	id, ok := ctx.Value(config.UserIDContextKey).(uuid.UUID)
	return id, ok
}

func setSessionCookie(w http.ResponseWriter, value string, expires time.Time) {
	if expires.IsZero() {
		slog.Error("setSessionCookie called with zero expires time, cookie not set")
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     config.SessionCookieName,
		Value:    value,
		Expires:  expires,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
}

func clearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     config.SessionCookieName,
		MaxAge:   -1,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
}

// JSON response helpers

type errorResponse struct {
	Status string `json:"status"`
	Data   string `json:"data"`
}

func writeJSONSuccess(w http.ResponseWriter, statusCode int, res any) {
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(res)
}

func writeJSONError(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(&errorResponse{
		Status: "error",
		Data:   message,
	})
}

func handleClientError(ctx context.Context, w http.ResponseWriter, err error, message string) {
	slog.WarnContext(ctx, "client error", "error", err, "message", message)
	writeJSONError(w, http.StatusBadRequest, message)
}

func handleServerError(ctx context.Context, w http.ResponseWriter, err error, message string) {
	slog.ErrorContext(ctx, "server error", "error", err, "message", message)
	writeJSONError(w, http.StatusInternalServerError, "internal server error")
}
