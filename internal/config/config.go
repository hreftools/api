package config

import "time"

const EmailVerificationTokenExpiryDuration = 24 * time.Hour
const PasswordResetTokenExpiryDuration = 1 * time.Hour

const SessionExpiryDuration = 30 * 24 * time.Hour
const SessionRenewalThreshold = 15 * 24 * time.Hour

type contextKey string

const UserIDContextKey contextKey = "userID"

const TokenTypeSession = "session"
const TokenTypeAPI = "token"

const SessionCookieName = "session_id"
