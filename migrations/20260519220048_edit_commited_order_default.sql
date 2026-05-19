-- +goose Up
ALTER TABLE comments
    ALTER COLUMN committed_order SET DEFAULT 1;
-- +goose Down
ALTER TABLE comments
    ALTER COLUMN committed_order DROP DEFAULT