-- +goose Up
CREATE INDEX index_guid_table ON storage(id);

-- +goose Down
DROP INDEX index_guid_table;
