-- name: User :one
SELECT * FROM users WHERE uuid = $1;

-- name: Users :many
SELECT * FROM users WHERE uuid != $1;

-- name: UserByEmail :one
SELECT * from users WHERE "email" = $1;

-- name: AvailableUsers :many
SELECT *
FROM users
    JOIN file_crypto_keys ON users.uuid = file_crypto_keys.user_uuid
WHERE users.uuid != $1 AND file_crypto_keys.file_uuid = $2;

-- name: UsersForShare :many
SELECT users.*
FROM users
    LEFT OUTER JOIN file_crypto_keys ON
    users.uuid = file_crypto_keys.user_uuid AND
    file_crypto_keys.file_uuid = $1
WHERE file_crypto_keys.user_uuid IS NULL;

-- name: CreateUser :one
INSERT INTO users (uuid, password_hash, email, name, public_key)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;