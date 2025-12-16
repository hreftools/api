package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"strings"

	"github.com/zapi-sh/api/internal/db"
	"github.com/zapi-sh/api/internal/response"
	"github.com/zapi-sh/api/internal/store"
	"github.com/zapi-sh/api/internal/utils"
)

type UserCreateBody struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

var reservedUsernames = map[string]bool{
	"admin": true,
}

func (b *UserCreateBody) Validate() error {
	if len(b.Username) < store.UserUsernameLengthMin {
		return errors.New("username must be min 3 characters")
	}

	if len(b.Username) > store.UserUsernameLengthMax {
		return errors.New("username must be max 32 characters")
	}

	if b.Username != strings.ToLower(b.Username) {
		return errors.New("username must be lowercase")
	}

	if !regexp.MustCompile(`^[a-z0-9_-]+$`).MatchString(b.Username) {
		return errors.New("username can only contain lowercase characters, numbers, hyphens, and underscores")
	}

	if strings.HasPrefix(b.Username, "-") || strings.HasPrefix(b.Username, "_") {
		return errors.New("username cannot start with hyphen or underscore")
	}

	if strings.HasSuffix(b.Username, "-") || strings.HasSuffix(b.Username, "_") {
		return errors.New("username cannot end with hyphen or underscore")
	}

	if len(b.Password) < store.UserPasswordLengthMin {
		return errors.New("password must be between 12")
	}

	return nil
}

type UserCreateResponse struct {
	Status string  `json:"status"`
	Data   db.User `json:"data"`
}

func UserCreate(store *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body UserCreateBody
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&body); err != nil {
			response.HandleClientError(w, err, "invalid request body")
			return
		}

		if err := body.Validate(); err != nil {
			response.HandleClientError(w, err, err.Error())
			return
		}

		passwordHash, err := utils.PasswordHash(body.Password)
		if err != nil {
			response.HandleClientError(w, err, "failed to hash password")
		}

		rr, err := store.Users.Create(r.Context(), body.Username, body.Email, passwordHash)
		if err != nil {
			response.HandleDbError(w, err)
			return
		}

		response := &UserCreateResponse{
			Status: "ok",
			Data:   rr,
		}

		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
