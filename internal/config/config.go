package config

import "time"

const SessionExpiryDuration = 30 * 24 * time.Hour
const SessionRenewalThreshold = 15 * 24 * time.Hour
const UserIDContextKey = "userID"
