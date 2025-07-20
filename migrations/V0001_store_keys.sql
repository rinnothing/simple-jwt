-- +goose Up
CREATE TABLE keys 
(
    access_key VARCHAR(64) NOT NULL,
    refresh_key VARCHAR(64) NOT NULL,
    refresh_hash_key VARCHAR(64) NOT NULL
);

-- +goose Down
DROP TABLE keys;
