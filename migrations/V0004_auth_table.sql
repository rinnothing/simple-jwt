-- +goose Up
CREATE TABLE auth
(
    id UUID PRIMARY KEY REFERENCES storage(id) ON DELETE CASCADE,
    refresh_hash BINARY(60) NOT NULL,
    user_agent TEXT NOT NULL,
    ip TEXT NOT NULL
);

-- +goose Down
DROP TABLE auth;
