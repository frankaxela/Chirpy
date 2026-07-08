-- Active: 1783508668124@@localhost@5432
-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email)
VALUES (
    gen_random_uuid(),
    now(),
    now(),
    $1
)
RETURNING *;

-- name: DeleteAllUsers :exec
DELETE FROM users;