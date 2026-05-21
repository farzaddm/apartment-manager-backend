-- +goose Up
-- +goose StatementBegin

CREATE TYPE ticket_accessability AS ENUM (
    'public',
    'private'
);

ALTER TABLE tickets
ADD COLUMN "accessability" ticket_accessability DEFAULT 'public';

UPDATE public.tickets SET accessability = 'public'
WHERE id = '11111111-bbbb-bbbb-bbbb-bbbbbbbbbbbb';

UPDATE public.tickets SET accessability = 'private'
WHERE id = '22222222-bbbb-bbbb-bbbb-bbbbbbbbbbbb';

UPDATE public.tickets SET accessability = 'public'
WHERE id = '33333333-bbbb-bbbb-bbbb-bbbbbbbbbbbb';

UPDATE public.tickets SET accessability = 'private'
WHERE id = '44444444-bbbb-bbbb-bbbb-bbbbbbbbbbbb';

UPDATE public.tickets SET accessability = 'public'
WHERE id = '55555555-bbbb-bbbb-bbbb-bbbbbbbbbbbb';

UPDATE public.tickets SET accessability = 'private'
WHERE id = '66666666-bbbb-bbbb-bbbb-bbbbbbbbbbbb';

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
ALTER TABLE tickets DROP COLUMN IF EXISTS "accessability";
DROP TYPE IF EXISTS ticket_accessability;

-- +goose StatementEnd
