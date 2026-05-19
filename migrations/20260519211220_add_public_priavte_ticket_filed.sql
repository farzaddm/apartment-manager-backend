-- +goose Up

CREATE TYPE ticket_accessability AS ENUM (
    'public',
    'private'
);
ALTER TABLE tickets ADD COLUMN "accessability" ticket_accessability DEFAULT 'public';

-- +goose Down
ALTER TABLE tickets DROP COLUMN IF EXISTS "accessability";
DROP TYPE IF EXISTS ticket_accessability;