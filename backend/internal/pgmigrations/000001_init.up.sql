CREATE TABLE users (
    uuid uuid PRIMARY KEY,
    password_hash text NOT NULL,
    email text NOT NULL UNIQUE,
    name text NOT NULL,
    public_key bytea NOT NULL
);

CREATE TABLE files (
    uuid uuid PRIMARY KEY,
    user_uuid uuid NOT NULL,
    name text NOT NULL,
    size bigint NOT NULL DEFAULT 0,
    is_crypt boolean NOT NULL DEFAULT false,
    created_at timestamp with time zone NOT NULL DEFAULT now(),
    FOREIGN KEY (user_uuid) REFERENCES users (uuid)
);

CREATE TABLE file_crypto_keys (
    file_uuid uuid NOT NULL,
    user_uuid uuid NOT NULL,
    symmetric_key text NOT NULL,
    PRIMARY KEY (user_uuid, file_uuid),
    FOREIGN KEY (user_uuid) REFERENCES users (uuid),
    FOREIGN KEY (file_uuid) REFERENCES files (uuid) ON DELETE CASCADE
);

CREATE UNIQUE INDEX "file_crypto_keys_file_uuid_user_uuid_key" ON file_crypto_keys (file_uuid, user_uuid);

