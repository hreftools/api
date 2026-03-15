package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hreftools/api/internal/response"
	"github.com/hreftools/api/internal/store"
	"github.com/hreftools/api/internal/utils"
	"github.com/hreftools/api/internal/validator"
)

type AuthResetPasswordConfirmBody struct {
	Token    string `json:"token"`
	Password string `json:"password"`
}

func (b *AuthResetPasswordConfirmBody) Normalize() {
	b.Token = strings.TrimSpace(b.Token)
}

func (b *AuthResetPasswordConfirmBody) Validate() error {
	if err := validator.Token(b.Token); err != nil {
		return err
	}

	if err := validator.Password(b.Password); err != nil {
		return err
	}

	return nil
}

type AuthResetPasswordConfirmResponse struct {
	Status string `json:"status"`
	Data   string `json:"data"`
}

func AuthResetPasswordConfirm(s *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body AuthResetPasswordConfirmBody
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&body); err != nil {
			response.HandleClientError(w, err, "invalid request body")
			return
		}

		body.Normalize()

		if err := body.Validate(); err != nil {
			response.HandleClientError(w, err, err.Error())
			return
		}

		token, _ := uuid.Parse(body.Token)

		u, err := s.Users.GetByPasswordResetToken(r.Context(), token)
		if err != nil {
			response.HandleDbError(w, err)
			return
		}

		if u.PasswordResetTokenExpiresAt != nil && u.PasswordResetTokenExpiresAt.Before(time.Now()) {
			response.HandleClientError(w, errors.New("token has expired"), "token has expired")
			return
		}

		passwordHash, err := utils.PasswordHash(body.Password)
		if err != nil {
			response.HandleServerError(w, err, "failed to hash password")
			return
		}

		_, err = s.Users.ResetPassword(r.Context(), u.ID, passwordHash)
		if err != nil {
			response.HandleDbError(w, err)
			return
		}

		go s.Tokens.DeleteAllByUserID(context.Background(), u.ID)

		res := AuthResetPasswordConfirmResponse{
			Status: "ok",
			Data:   "ok",
		}

		response.WriteJSONSuccess(w, http.StatusOK, res)
	}
}
