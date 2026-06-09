-- +goose Up
ALTER TABLE tickets
ADD COLUMN apartment_id UUID;

ALTER TABLE tickets
ADD CONSTRAINT fk_tickets_apartment
FOREIGN KEY (apartment_id) REFERENCES apartments(id)
ON UPDATE CASCADE
ON DELETE CASCADE;

CREATE INDEX idx_tickets_apartment_id ON tickets(apartment_id);

UPDATE tickets
SET apartment_id = 'a0000000-0000-0000-0000-000000000000'
WHERE id IN (
  '11111111-bbbb-bbbb-bbbb-bbbbbbbbbbbb',
  '22222222-bbbb-bbbb-bbbb-bbbbbbbbbbbb',
  '33333333-bbbb-bbbb-bbbb-bbbbbbbbbbbb',
  '44444444-bbbb-bbbb-bbbb-bbbbbbbbbbbb',
  '55555555-bbbb-bbbb-bbbb-bbbbbbbbbbbb',
  '66666666-bbbb-bbbb-bbbb-bbbbbbbbbbbb',
  '77777777-bbbb-bbbb-bbbb-bbbbbbbbbbbb',
  '88888888-bbbb-bbbb-bbbb-bbbbbbbbbbbb',
  '99999999-bbbb-bbbb-bbbb-bbbbbbbbbbbb',
  'aaaaaaaa-bbbb-bbbb-bbbb-bbbbbbbbbbbb',
  'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb'
);

ALTER TABLE tickets
ALTER COLUMN apartment_id SET NOT NULL;

-- +goose Down
DROP INDEX IF EXISTS idx_tickets_apartment_id;

ALTER TABLE tickets
DROP CONSTRAINT IF EXISTS fk_tickets_apartment;

ALTER TABLE tickets
DROP COLUMN IF EXISTS apartment_id;
