-- name: CreateChirp :one
INSERT INTO chirps (id, body, user_id, created_at, updated_at)
VALUES (
    gen_random_uuid(),
    $1,
    $2,
    $3,
    $4
)
RETURNING *;

-- name: GetChirps :many
SELECT * FROM chirps
ORDER BY created_at ASC;

-- name: GetChirpsByAuthor :many
SELECT * FROM chirps
WHERE user_id = $1;

-- name: GetChirp :one
SELECT * FROM chirps
WHERE id = $1;

-- name: DeleteAllChirps :exec
DELETE FROM chirps;

-- name: DeleteChirp :exec
DELETE FROM chirps
WHERE id = $1;