-- name: CreateUser :one
INSERT INTO users (email, password_hash, first_name, last_name)
VALUES ($1, $2, $3, $4)
RETURNING id, email, first_name, last_name, created_at, updated_at;


-- name: UpdateGender :exec
UPDATE users
SET gender = $1, updated_at = now()
WHERE id = $2;


-- name: UpdatePhoneNumber :exec
UPDATE users
SET phone_number = $1, updated_at = now()
WHERE id = $2;


-- name: UpdateUserVerified :exec
UPDATE users
SET is_verified = $1, updated_at = now()
WHERE id = $2;

-- name: UpdateUserPassword :exec
UPDATE users
SET password_hash = $1, updated_at = now()
WHERE id = $2;

-- name: GetUserByEmail :one
SELECT id, email, password_hash, is_verified, first_name, last_name
FROM users
WHERE email = $1;


-- name: GetUserByID :one
SELECT id, email, is_verified, first_name, last_name
FROM users
WHERE id = $1;



