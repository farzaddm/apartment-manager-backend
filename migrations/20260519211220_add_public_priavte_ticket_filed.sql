-- +goose Up
-- +goose StatementBegin

CREATE TYPE ticket_accessibility AS ENUM (
    'public',
    'private'
);

ALTER TABLE tickets
ADD COLUMN "accessibility" ticket_accessibility DEFAULT 'public';
UPDATE public.tickets SET accessibility = 'public'
WHERE id = '11111111-bbbb-bbbb-bbbb-bbbbbbbbbbbb';

UPDATE public.tickets SET accessibility = 'private'
WHERE id = '22222222-bbbb-bbbb-bbbb-bbbbbbbbbbbb';

UPDATE public.tickets SET accessibility = 'public'
WHERE id = '33333333-bbbb-bbbb-bbbb-bbbbbbbbbbbb';

UPDATE public.tickets SET accessibility = 'private'
WHERE id = '44444444-bbbb-bbbb-bbbb-bbbbbbbbbbbb';

UPDATE public.tickets SET accessibility = 'public'
WHERE id = '55555555-bbbb-bbbb-bbbb-bbbbbbbbbbbb';

UPDATE public.tickets SET accessibility = 'private'
WHERE id = '66666666-bbbb-bbbb-bbbb-bbbbbbbbbbbb';

UPDATE public.tickets SET accessibility = 'public'
WHERE id = '77777777-bbbb-bbbb-bbbb-bbbbbbbbbbbb';

UPDATE public.tickets SET accessibility = 'private'
WHERE id = '88888888-bbbb-bbbb-bbbb-bbbbbbbbbbbb';

UPDATE public.tickets SET accessibility = 'public'
WHERE id = '99999999-bbbb-bbbb-bbbb-bbbbbbbbbbbb';

UPDATE public.tickets SET accessibility = 'private'
WHERE id = 'aaaaaaaa-bbbb-bbbb-bbbb-bbbbbbbbbbbb';

UPDATE public.tickets SET accessibility = 'public'
WHERE id = 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb';

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
ALTER TABLE tickets DROP COLUMN IF EXISTS "accessibility";
DROP TYPE IF EXISTS ticket_accessibility;

-- +goose StatementEnd
