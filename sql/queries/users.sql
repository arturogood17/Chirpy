-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;

-- name: ResetUsersTable :exec
DELETE FROM users;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: UpdateUser :one
UPDATE users
SET updated_at = NOW(), hashed_password = $1, email = $2
WHERE id = $3
RETURNING *;

-- name: ChirpyRed :one
UPDATE users
SET is_chirpy_red = true
WHERE id = $1
RETURNING *;
