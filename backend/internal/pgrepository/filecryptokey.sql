-- name: UpsertFileCryptoKey :one
INSERT INTO file_crypto_keys (file_uuid, user_uuid, symmetric_key)
VALUES ($1, $2, $3)
ON CONFLICT (file_uuid, user_uuid) DO
UPDATE SET symmetric_key = EXCLUDED.symmetric_key
RETURNING *;

-- name: DeleteFileAccess :exec
DELETE FROM file_crypto_keys
WHERE file_crypto_keys.user_uuid = @recipient_uuid
  AND file_crypto_keys.file_uuid = @file_uuid
  AND file_crypto_keys.file_uuid = (
    SELECT files.uuid FROM files
    WHERE files.user_uuid = @owner_uuid
    );

-- name: DeleteAlleFileAccess :exec
DELETE FROM file_crypto_keys
WHERE file_crypto_keys.user_uuid = $1;