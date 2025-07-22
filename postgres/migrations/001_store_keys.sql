-- +goose Up
CREATE TABLE keys 
(
    access_key BYTEA NOT NULL,
    refresh_key BYTEA NOT NULL,
    refresh_hash_key BYTEA NOT NULL,
    created_at TIMESTAMP DEFAULT now() NOT NULL
);

-- +goose Down
DROP TABLE keys;
