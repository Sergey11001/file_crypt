-- name: File :one
SELECT sqlc.embed(files),
       file_crypto_keys.symmetric_key
FROM files
LEFT JOIN file_crypto_keys ON files.uuid = file_crypto_keys.file_uuid
WHERE file_crypto_keys.user_uuid = $1 AND files.uuid = $2;

-- name: FileByUUID :one
SELECT *
FROM files
WHERE files.uuid = $1;

-- name: CommonFile :one
SELECT sqlc.embed(files),
       file_crypto_keys.symmetric_key
FROM files
         LEFT JOIN file_crypto_keys ON files.uuid = file_crypto_keys.file_uuid
WHERE file_crypto_keys.user_uuid IS NULL AND files.uuid = $1;

-- name: Files :many
SELECT sqlc.embed(files),
       file_crypto_keys.symmetric_key
FROM files
LEFT JOIN file_crypto_keys ON files.uuid = file_crypto_keys.file_uuid AND files.user_uuid = file_crypto_keys.user_uuid
WHERE files.user_uuid = $1;

-- name: AvailableFiles :many
SELECT sqlc.embed(files),
       file_crypto_keys.symmetric_key
FROM files
         JOIN file_crypto_keys ON files.uuid = file_crypto_keys.file_uuid
WHERE file_crypto_keys.user_uuid = $1 AND files.user_uuid != file_crypto_keys.user_uuid;

-- name: CreateFile :one
INSERT INTO files (uuid, user_uuid, name, size, is_crypt)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: DeleteFile :exec
DELETE FROM files
WHERE user_uuid = $1 AND uuid = $2;

-- name: DeleteAllUserFiles :many
DELETE FROM files
WHERE user_uuid = $1
RETURNING files.uuid;