-- +goose Up
CREATE TYPE announcement_order AS ENUM (
    'warning',
    'very_important',
    'important',
    'other'
);


ALTER TABLE announcements ADD COLUMN "order" announcement_order DEFAULT 'other';
ALTER TABLE announcements ADD COLUMN is_pinned BOOLEAN NOT NULL DEFAULT false;

UPDATE announcements SET "order" = 'warning' WHERE id IN ('11111111-aaaa-aaaa-aaaa-aaaaaaaaaaaa', '66666666-aaaa-aaaa-aaaa-aaaaaaaaaaaa', '10101010-aaaa-aaaa-aaaa-aaaaaaaaaaaa', '14141414-aaaa-aaaa-aaaa-aaaaaaaaaaaa', '17171717-aaaa-aaaa-aaaa-aaaaaaaaaaaa');
UPDATE announcements SET "order" = 'very_important' WHERE id IN ('44444444-aaaa-aaaa-aaaa-aaaaaaaaaaaa', '77777777-aaaa-aaaa-aaaa-aaaaaaaaaaaa', '20202020-aaaa-aaaa-aaaa-aaaaaaaaaaaa');
UPDATE announcements SET "order" = 'important' WHERE id IN ('55555555-aaaa-aaaa-aaaa-aaaaaaaaaaaa', '88888888-aaaa-aaaa-aaaa-aaaaaaaaaaaa', '99999999-aaaa-aaaa-aaaa-aaaaaaaaaaaa', '12121212-aaaa-aaaa-aaaa-aaaaaaaaaaaa', '16161616-aaaa-aaaa-aaaa-aaaaaaaaaaaa');

UPDATE announcements SET is_pinned = true WHERE id IN ('11111111-aaaa-aaaa-aaaa-aaaaaaaaaaaa', '44444444-aaaa-aaaa-aaaa-aaaaaaaaaaaa', '66666666-aaaa-aaaa-aaaa-aaaaaaaaaaaa', '77777777-aaaa-aaaa-aaaa-aaaaaaaaaaaa', '14141414-aaaa-aaaa-aaaa-aaaaaaaaaaaa', '16161616-aaaa-aaaa-aaaa-aaaaaaaaaaaa');


-- +goose Down
ALTER TABLE announcements DROP COLUMN IF EXISTS "order";
ALTER TABLE announcements DROP COLUMN IF EXISTS is_pinned;
DROP TYPE IF EXISTS announcement_order;