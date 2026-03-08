# Plan: Session-based login (authSignin) handler

## Context

The API has signup, email verification, and resend verification handlers. The next step is login. The user wants session-based auth with HttpOnly cookies. API tokens for third-party/plugin use will be a separate feature later.

**Key decisions:**
- Cookie-only — no token in response body
- Require verified email to log in
- Sliding expiration: 30 days of inactivity → session expires. Refresh is throttled (future middleware concern, not part of this handler)
- Generic error message ("invalid email or password") for both wrong email and wrong password

## Steps

### 1. `sql/migrations/000001_init.up.sql` (modify)
Append the sessions table and trigger after the resources section:
```sql
-- sessions
CREATE TABLE sessions (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TRIGGER update_sessions_updated_at
BEFORE UPDATE ON sessions
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
```

### 2. `sql/migrations/000001_init.down.sql` (modify)
Add drops for sessions trigger and table (before the existing drops, since sessions references users).

### 3. `sql/queries/sessions.sql` (create)
Queries following the formatting style in `sql/queries/users.sql`:
- `CreateSession :one` — INSERT user_id + expires_at, RETURNING *
- `GetSessionById :one` — SELECT by id
- `UpdateSessionExpiresAt :one` — UPDATE expires_at WHERE id, RETURNING * (for future sliding expiration middleware)
- `DeleteSession :exec` — DELETE by id (for future signout)

### 4. `sqlc.yml` (modify)
Add overrides for sessions timestamp columns (non-pointer `time.Time`, matching resources/users pattern):
- `sessions.expires_at`
- `sessions.created_at`
- `sessions.updated_at`

### 5. `internal/store/store_sessions.go` (create)
Following the exact pattern of `internal/store/store_users.go`:
- `SessionCreateParams` and `SessionUpdateExpiresAtParams` structs
- `SessionStore` interface with `Create`, `GetById`, `UpdateExpiresAt`, `Delete`
- `sessionStore` concrete type backed by `db.Querier`
- `NewSessionStore(queries db.Querier) SessionStore` constructor

### 6. `internal/store/store.go` (modify)
Add `Sessions SessionStore` field to `Store` struct (alphabetically between Resources and Users) and wire `NewSessionStore(queries)` in `New()`.

### 7. `internal/handlers/constants.go` (modify)
Add two constants:
- `SessionExpiryDuration = 30 * 24 * time.Hour`
- `SessionCookieName = "session_id"`

### 8. `internal/handlers/authSignin.go` (create)
Handler following the established pattern (decode → normalize → validate → business logic → response):
- `AuthSigninBody` with `Email` and `Password` fields
- `Normalize()` — lowercase + trim email, leave password untouched
- `Validate()` — reuse `validator.Email()` and `validator.Password()`
- `AuthSigninResponse` with `Status string` and `Data string`
- Handler flow:
  1. Decode + DisallowUnknownFields
  2. Normalize + Validate
  3. `s.Users.GetByEmail()` → on error, return 401 "invalid email or password"
  4. `utils.PasswordValidate()` → on failure, return 401 "invalid email or password"
  5. Check `u.EmailVerified` → if false, return 403 "email not verified"
  6. `s.Sessions.Create()` with expires_at = now + 30 days
  7. `http.SetCookie()` — HttpOnly, Secure, SameSite=Lax, Path="/", name="session_id"
  8. Respond 200 `{"status":"ok","data":"ok"}`

### 9. `internal/handlers/authSignin_test.go` (create)
Integration tests matching the pattern in `authSignup_test.go`:
- "fails on incorrect body" → 400
- "fails on unexpected field in body" → 400
- "fails on invalid request body" (empty email) → 400
- "fails on non-existent email" → 401 "invalid email or password"
- "fails on wrong password" → 401 "invalid email or password"
- "fails on unverified email" → 403 "email not verified"
- "success" → 200, check response body + Set-Cookie header (name, HttpOnly, Secure, valid UUID value, Expires ~30 days)

Test users seeded via `s.Users.Create()` with pre-hashed passwords using `utils.PasswordHash()`.

### 10. `internal/handlers/helpers_test.go` (modify)
Add `sessions` to the TRUNCATE statement: `"TRUNCATE users, resources, sessions RESTART IDENTITY CASCADE"`

### 11. `internal/server/server.go` (modify)
Uncomment and wire the signin route:
```go
mux.HandleFunc("POST /auth/signin", handlers.AuthSignin(s))
```

## Files auto-generated (by sqlc generate, not manually edited)
- `internal/db/models.go` — gains `Session` struct
- `internal/db/querier.go` — gains session methods
- `internal/db/sessions.sql.go` — new file

## Verification

After implementation, the user should run:
1. Re-run the migration (drop + up) on dev and test databases since the init migration was modified
2. `sqlc generate`
3. `go build ./...`
4. `make test`
