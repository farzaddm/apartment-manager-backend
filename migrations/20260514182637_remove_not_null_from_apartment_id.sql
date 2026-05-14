-- +goose Up
ALTER TABLE users ALTER COLUMN apartment_id DROP NOT NULL;

-- +goose Down
ALTER TABLE users ALTER COLUMN apartment_id SET NOT NULL;
