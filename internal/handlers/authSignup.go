package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/mail"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hreftools/api/internal/db"
	"github.com/hreftools/api/internal/emails"
	"github.com/hreftools/api/internal/response"
	"github.com/hreftools/api/internal/store"
	"github.com/hreftools/api/internal/utils"
	"github.com/resend/resend-go/v3"
)

type AuthSignupBody struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (b *AuthSignupBody) Normalize() {
	b.Username = strings.ToLower(strings.TrimSpace(b.Username))
	b.Email = strings.ToLower(strings.TrimSpace(b.Email))
}

func (b *AuthSignupBody) Validate() error {
	// username
	if len(b.Username) == 0 {
		return errors.New("username is required")
	}

	if len(b.Username) < store.UserUsernameLengthMin {
		return errors.New("username must be min 3 characters")
	}

	if len(b.Username) > store.UserUsernameLengthMax {
		return errors.New("username must be max 32 characters")
	}

	if !store.UserPattern.MatchString(b.Username) {
		return errors.New("username can only contain lowercase characters, numbers, hyphens, and underscores")
	}

	if strings.HasPrefix(b.Username, "-") || strings.HasPrefix(b.Username, "_") {
		return errors.New("username cannot start with hyphen or underscore")
	}

	if strings.HasSuffix(b.Username, "-") || strings.HasSuffix(b.Username, "_") {
		return errors.New("username cannot end with hyphen or underscore")
	}

	if reserved := reservedUsernames[b.Username]; reserved {
		return errors.New("username is reserved")
	}

	// email
	if len(b.Email) == 0 {
		return errors.New("email is required")
	}

	// validate format RFC 5322
	if _, err := mail.ParseAddress(b.Email); err != nil {
		return errors.New("email format is invalid")
	}

	// limit length as per smtp spec RFC 5321
	if len(b.Email) > 254 {
		return errors.New("email must be at most 254 characters")
	}

	// password
	if len(b.Password) == 0 {
		return errors.New("password is required")
	}

	if len(b.Password) < store.UserPasswordLengthMin {
		return errors.New("password must be at least 12 characters")
	}

	return nil
}

type AuthSignupResponse struct {
	Status string  `json:"status"`
	Data   db.User `json:"data"`
}

func AuthSignup(s *store.Store, emailSender emails.EmailSender) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body AuthSignupBody
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

		passwordHash, err := utils.PasswordHash(body.Password)
		if err != nil {
			response.HandleClientError(w, err, "failed to hash password")
			return
		}

		email := body.Email
		username := body.Username
		emailVerified := false
		emailVerificationToken := uuid.NullUUID{Valid: true, UUID: uuid.New()}
		emailVerificationTokenExpiresAt := time.Now().Add(store.UserVerificationTokenExpiration)
		isAdmin := false
		isPro := false
		u, err := s.Users.Create(r.Context(), email, emailVerified, emailVerificationToken, &emailVerificationTokenExpiresAt, passwordHash, username, isAdmin, isPro)
		if err != nil {
			response.HandleDbError(w, err)
			return
		}

		emailVerifyData := emails.EmailVerifyData{
			Username:  username,
			Email:     email,
			Token:     emailVerificationToken.UUID.String(),
			ExpiresAt: emailVerificationTokenExpiresAt.Format(time.RFC1123),
		}
		bodyHtml, err := emails.EmailVerifyRenderHtml(emailVerifyData)
		if err != nil {
			response.HandleClientError(w, err, "failed to render html email template")
			return
		}
		bodyText, err := emails.EmailVerifyRenderTxt(emailVerifyData)
		if err != nil {
			response.HandleClientError(w, err, "failed to render text email template")
			return
		}
		params := &resend.SendEmailRequest{
			From:    "href.tools <auth@mail.href.tools>",
			To:      []string{email},
			Text:    bodyText,
			Html:    bodyHtml,
			Subject: "Hello from href.tools",
			ReplyTo: "auth@mail.href.tools",
		}

		err = emailSender.Send(params)
		if err != nil {
			response.HandleClientError(w, err, err.Error())
			return
		}

		response := &AuthSignupResponse{
			Status: "ok",
			Data:   u,
		}

		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
