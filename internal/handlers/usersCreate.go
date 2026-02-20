package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/hreftools/api/internal/db"
	"github.com/hreftools/api/internal/response"
	"github.com/hreftools/api/internal/store"
	"github.com/hreftools/api/internal/utils"
	"github.com/hreftools/api/internal/validator"
)

type UserCreateBody struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	IsAdmin  *bool  `json:"isAdmin"`
	IsPro    *bool  `json:"isPro"`
}

func (b *UserCreateBody) Validate() error {
	if err := validator.Username(b.Username); err != nil {
		return err
	}

	if err := validator.Email(b.Email); err != nil {
		return err
	}

	if err := validator.Password(b.Password); err != nil {
		return err
	}

	if b.IsAdmin == nil {
		return errors.New("isAdmin field is required")
	}

	if b.IsPro == nil {
		return errors.New("isPro field is required")
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
			return
		}
		u, err := store.Users.Create(r.Context(), strings.TrimSpace(body.Email), true, uuid.NullUUID{}, nil, passwordHash, strings.TrimSpace(body.Username), *body.IsAdmin, *body.IsPro)
		if err != nil {
			response.HandleDbError(w, err)
			return
		}

		response := &UserCreateResponse{
			Status: "ok",
			Data:   u,
		}

		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
