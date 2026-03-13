package config

import "time"

const TokenExpiryDuration = 24 * time.Hour

const SessionExpiryDuration = 30 * 24 * time.Hour
const SessionRenewalThreshold = 15 * 24 * time.Hour
const UserIDContextKey = "userID"

const TokenTypeSession = "session"
const TokenTypeAPI = "token"

const SessionCookieName = "session_id"
