-- +goose Up
INSERT INTO apartments (id, name, province, city, address, postal_code)
VALUES (
           'c0000000-0000-0000-0000-000000000000',
           'Grand Horizon Luxury Apartments',
           'Ontario',
           'Toronto',
           '123 Skyview Avenue, Suite 100',
           'M5V 2T6'
       );

INSERT INTO units (id, apartment_id, user_id, unit_number, floor)
VALUES
    ('11111113-eeee-eeee-eeee-eeeeeeeeeeee', 'c0000000-0000-0000-0000-000000000000', NULL, '101', 1),
    ('22222224-eeee-eeee-eeee-eeeeeeeeeeee', 'c0000000-0000-0000-0000-000000000000', NULL, '102', 1),
    ('33333335-eeee-eeee-eeee-eeeeeeeeeeee', 'c0000000-0000-0000-0000-000000000000', NULL, '201', 2),
    ('44444446-eeee-eeee-eeee-eeeeeeeeeeee', 'c0000000-0000-0000-0000-000000000000', NULL, '202', 2),
    ('55555557-eeee-eeee-eeee-eeeeeeeeeeee', 'c0000000-0000-0000-0000-000000000000', NULL, '301', 3),
    ('66666668-eeee-eeee-eeee-eeeeeeeeeeee', 'c0000000-0000-0000-0000-000000000000', NULL, '302', 3);

INSERT INTO invite_codes (id, apartment_id, unit_id, code, expires_at)
VALUES
    ('11111112-9999-9999-9999-999999999991', 'c0000000-0000-0000-0000-000000000000', '55555557-eeee-eeee-eeee-eeeeeeeeeeee', 'INV-30D-WELCOME', NOW() + INTERVAL '30 days'),
    ('11111112-9999-9999-9999-999999999992', 'c0000000-0000-0000-0000-000000000000', '66666668-eeee-eeee-eeee-eeeeeeeeeeee', 'INV-3Q2-WELCOME', NOW() + INTERVAL '30 days'),
    ('11111112-9999-9999-9999-999999999993', 'c0000000-0000-0000-0000-000000000000', '11111113-eeee-eeee-eeee-eeeeeeeeeeee', 'INV-ERPIRED-99', NOW() - INTERVAL '12 day'),
    ('11111112-9999-9999-9999-999999999994', 'c0000000-0000-0000-0000-000000000000', '22222224-eeee-eeee-eeee-eeeeeeeeeeee', 'CODE-RLPHA-444', NOW() + INTERVAL '10 days'),
    ('11111112-9999-9999-9999-999999999995', 'c0000000-0000-0000-0000-000000000000', '33333335-eeee-eeee-eeee-eeeeeeeeeeee', 'CODE-BQTA-555',  NOW() + INTERVAL '15 days'),
    ('11111112-9999-9999-9999-999999999996', 'c0000000-0000-0000-0000-000000000000', '44444446-eeee-eeee-eeee-eeeeeeeeeeee', 'CODE-TTMMA-666', NOW() + INTERVAL '20 days');

-- +goose Down
SELECT 'down SQL query';
