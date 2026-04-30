package config

import (
	"fmt"
	"os"
	"strings"
	"time"
)

const EmailVerificationTokenExpiryDuration = 24 * time.Hour
const PasswordResetTokenExpiryDuration = 1 * time.Hour

const SessionExpiryDuration = 30 * 24 * time.Hour
const SessionRenewalThreshold = 15 * 24 * time.Hour

type contextKey string

const UserIDContextKey contextKey = "userID"

const SessionCookieName = "session_id"

type Config struct {
	Port         string
	DatabaseURL  string
	ResendAPIKey string
	AppURL       string
}

func LoadConfig() (*Config, error) {
	cfg := &Config{
		Port:         os.Getenv("PORT"),
		DatabaseURL:  os.Getenv("DATABASE_URL"),
		ResendAPIKey: os.Getenv("RESEND_API_KEY"),
		AppURL:       strings.TrimSuffix(os.Getenv("APP_URL"), "/"),
	}

	var missing []string

	if cfg.Port == "" {
		missing = append(missing, "PORT")
	}
	if cfg.DatabaseURL == "" {
		missing = append(missing, "DATABASE_URL")
	}
	if cfg.ResendAPIKey == "" {
		missing = append(missing, "RESEND_API_KEY")
	}
	if cfg.AppURL == "" {
		missing = append(missing, "APP_URL")
	}

	if len(missing) > 0 {
		return nil, fmt.Errorf("missing required environment variables: %s", strings.Join(missing, ", "))
	}

	return cfg, nil
}
