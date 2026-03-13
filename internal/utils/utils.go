package utils

import (
	"context"

	"github.com/google/uuid"
	"github.com/hreftools/api/internal/config"
	"golang.org/x/crypto/bcrypt"
)

func UserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	id, ok := ctx.Value(config.UserIDContextKey).(uuid.UUID)
	return id, ok
}

func PasswordHash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func PasswordValidate(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
