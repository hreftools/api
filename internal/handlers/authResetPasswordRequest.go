package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hreftools/api/internal/config"
	"github.com/hreftools/api/internal/emails"
	"github.com/hreftools/api/internal/response"
	"github.com/hreftools/api/internal/store"
	"github.com/hreftools/api/internal/validator"
)

type AuthResetPasswordRequestBody struct {
	Email string `json:"email"`
}

func (b *AuthResetPasswordRequestBody) Normalize() {
	b.Email = strings.ToLower(strings.TrimSpace(b.Email))
}

func (b *AuthResetPasswordRequestBody) Validate() error {
	if err := validator.Email(b.Email); err != nil {
		return err
	}

	return nil
}

type AuthResetPasswordRequestResponse struct {
	Status string `json:"status"`
	Data   string `json:"data"`
}

func AuthResetPasswordRequest(s *store.Store, emailSender emails.EmailSender) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body AuthResetPasswordRequestBody
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

		u, err := s.Users.GetByEmail(r.Context(), body.Email)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				response.WriteJSONSuccess(w, http.StatusOK, &AuthResetPasswordRequestResponse{
					Status: "ok",
					Data:   "ok",
				})
				return
			}
			response.HandleServerError(w, err, "failed to query user")
			return
		}

		// rate limit: if existing token age < 5 minutes, return 429
		if u.PasswordResetTokenExpiresAt != nil {
			tokenAge := config.PasswordResetTokenExpiryDuration - time.Until(*u.PasswordResetTokenExpiresAt)
			if tokenAge < time.Minute*5 {
				log.Println("password reset email already sent, please wait before requesting a new one")
				response.WriteJSONError(w, http.StatusTooManyRequests, "password reset email already sent, please wait before requesting a new one")
				return
			}
		}

		token := uuid.NullUUID{Valid: true, UUID: uuid.New()}

		templateParams := emails.AuthResetPasswordRequestParams{
			Token: token.UUID.String(),
		}
		bodyHtml, err := emails.RenderTemplateHtml(emails.AuthResetPasswordRequestTemplateHtml, templateParams)
		if err != nil {
			response.HandleServerError(w, err, "failed to render html email template")
			return
		}
		bodyText, err := emails.RenderTemplateTxt(emails.AuthResetPasswordRequestTemplateTxt, templateParams)
		if err != nil {
			response.HandleServerError(w, err, "failed to render text email template")
			return
		}

		params := store.UserUpdatePasswordResetTokenParams{
			Id:                          u.ID,
			PasswordResetToken:          token,
			PasswordResetTokenExpiresAt: new(time.Now().Add(config.PasswordResetTokenExpiryDuration)),
		}
		_, err = s.Users.UpdatePasswordResetToken(r.Context(), params)
		if err != nil {
			response.HandleDbError(w, err)
			return
		}

		emailParams := emails.EmailSendParams{
			To:      []string{body.Email},
			Text:    bodyText,
			Html:    bodyHtml,
			Subject: "Password reset has been requested",
		}

		err = emailSender.Send(emailParams)
		if err != nil {
			log.Printf("Failed to send email: %v", err)
		}

		res := AuthResetPasswordRequestResponse{
			Status: "ok",
			Data:   "ok",
		}

		response.WriteJSONSuccess(w, http.StatusOK, res)
	}
}
