-- +goose Up
CREATE TYPE announcement_order AS ENUM (
    'warning',
    'very_important',
    'important',
    'other'
);


ALTER TABLE announcements ADD COLUMN "order" announcement_order DEFAULT 'other';
ALTER TABLE announcements ADD COLUMN is_pinned BOOLEAN NOT NULL DEFAULT false;



-- +goose Down
ALTER TABLE announcements DROP COLUMN IF EXISTS "order";
ALTER TABLE announcements DROP COLUMN IF EXISTS is_pinned;
DROP TYPE IF EXISTS announcement_order;