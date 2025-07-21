-- +goose Up
CREATE INDEX index_auth_table ON auth(id);

-- +goose Down
DROP INDEX index_auth_table;
