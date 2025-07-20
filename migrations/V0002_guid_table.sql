-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE storage
(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    guid TEXT NOT NULL
);

-- +goose Down
DROP TABLE storage;
