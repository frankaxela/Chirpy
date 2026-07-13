
-- name: CreateRefreshToken :one
-- Insert a new refresh token and return the created row
INSERT INTO refresh_tokens (
	user_id,
	token,
	created_at,
	updated_at,
    expires_at
)
VALUES ($1, $2, $3, $4, $5)
RETURNING user_id, token, created_at, updated_at, expires_at;

-- name: GetRefreshToken :one
-- Look up a refresh token by its token value
SELECT token, created_at, updated_at, user_id, expires_at, revoked_at
FROM refresh_tokens
WHERE token = $1;

-- name: RevokeRefreshToken :exec
-- Expire a refresh token by setting its revoked_at timestamp
UPDATE refresh_tokens
SET revoked_at = NOW(),
    updated_at = NOW()
WHERE token = $1;


