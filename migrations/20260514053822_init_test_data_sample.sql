-- +goose Up

-- ==========================================================
-- 1. APARTMENT
-- ==========================================================
INSERT INTO apartments (id, name, province, city, address, postal_code, unit_count)
VALUES ('7f23c713-379e-4b9d-8386-8d197600ef58', 'The Grand Residency', 'Ontario', 'Toronto', '777 Luxury Lane', 'M5V 3L9', 120);

-- ==========================================================
-- 2. USERS (Corrected Hexadecimal UUIDs)
-- ==========================================================
INSERT INTO users (id, apartment_id, first_name, last_name, username, email, phone, password, role, gender)
VALUES 
('00000000-0000-0000-0000-000000000001', '7f23c713-379e-4b9d-8386-8d197600ef58', 'Admin', 'User', 'admin', 'admin@test.com', '555-0001', 'hash', 'admin', 'male'),
('00000000-0000-0000-0000-000000000002', '7f23c713-379e-4b9d-8386-8d197600ef58', 'Manager', 'One', 'manager1', 'm1@test.com', '555-0002', 'hash', 'manager', 'female'),
('00000000-0000-0000-0000-000000000003', '7f23c713-379e-4b9d-8386-8d197600ef58', 'Alice', 'Resident', 'alice_res', 'alice@test.com', '555-1001', 'hash', 'resident', 'female'),
('00000000-0000-0000-0000-000000000004', '7f23c713-379e-4b9d-8386-8d197600ef58', 'Bob', 'Resident', 'bob_res', 'bob@test.com', '555-1002', 'hash', 'resident', 'male'),
('00000000-0000-0000-0000-000000000005', '7f23c713-379e-4b9d-8386-8d197600ef58', 'Charlie', 'Resident', 'charlie_res', 'charlie@test.com', '555-1003', 'hash', 'resident', 'male'),
('00000000-0000-0000-0000-000000000006', '7f23c713-379e-4b9d-8386-8d197600ef58', 'David', 'Resident', 'david_res', 'david@test.com', '555-1004', 'hash', 'resident', 'male'),
('00000000-0000-0000-0000-000000000007', '7f23c713-379e-4b9d-8386-8d197600ef58', 'Eve', 'Resident', 'eve_res', 'eve@test.com', '555-1005', 'hash', 'resident', 'female');

-- ==========================================================
-- 3. TAGS
-- ==========================================================
INSERT INTO tags (id, name)
VALUES 
('a0000000-0000-0000-0000-000000000001', 'Urgent'),
('a0000000-0000-0000-0000-000000000002', 'Maintenance'),
('a0000000-0000-0000-0000-000000000003', 'Plumbing'),
('a0000000-0000-0000-0000-000000000004', 'Security');

-- ==========================================================
-- 4. TICKETS
-- ==========================================================
INSERT INTO tickets (id, user_id, title, description, body, category, status)
VALUES 
('b0000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000003', 'Leaking Sink', 'Kitchen sink leak', 'Water is dripping slowly.', 'Maintenance', 'open'),
('b0000000-0000-0000-0000-000000000002', '00000000-0000-0000-0000-000000000004', 'No Hot Water', 'Bathroom shower', 'The water is freezing cold.', 'Maintenance', 'open'),
('b0000000-0000-0000-0000-000000000003', '00000000-0000-0000-0000-000000000005', 'Lost Key', 'Front door key lost', 'I lost my main entrance key.', 'Security', 'closed');

-- ==========================================================
-- 5. COMMENTS
-- ==========================================================
INSERT INTO comments (id, user_id, ticket_id, body, commited_order)
VALUES 
('c0000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000002', 'b0000000-0000-0000-0000-000000000001', 'Plumber scheduled.', 1),
('c0000000-0000-0000-0000-000000000002', '00000000-0000-0000-0000-000000000003', 'b0000000-0000-0000-0000-000000000001', 'Thanks.', 2);

-- ==========================================================
-- 6. ANNOUNCEMENTS
-- ==========================================================
INSERT INTO announcements (id, apartment_id, title, description, body, expired_date)
VALUES 
('d0000000-0000-0000-0000-000000000001', '7f23c713-379e-4b9d-8386-8d197600ef58', 'Roof Party', 'Summer celebration', 'Roof party this Friday.', '2026-12-31 23:59:59'),
('d0000000-0000-0000-0000-000000000002', '7f23c713-379e-4b9d-8386-8d197600ef58', 'Fire Drill', 'Safety inspection', 'Fire drill Monday morning.', '2026-12-31 23:59:59');

-- ==========================================================
-- 7. RULES & RULE ITEMS
-- ==========================================================
INSERT INTO rules (id, apartment_id, category)
VALUES ('e0000000-0000-0000-0000-000000000001', '7f23c713-379e-4b9d-8386-8d197600ef58', 'Pool Rules');

INSERT INTO rule_items (id, rule_id, body)
VALUES 
('f0000000-0000-0000-0000-000000000001', 'e0000000-0000-0000-0000-000000000001', 'No running on deck.'),
('f0000000-0000-0000-0000-000000000002', 'e0000000-0000-0000-0000-000000000001', 'No glass allowed.');

-- ==========================================================
-- 8. INVITE CODES
-- ==========================================================
INSERT INTO invite_codes (id, apartment_id, code, expires_at)
VALUES ('90000000-0000-0000-0000-000000000001', '7f23c713-379e-4b9d-8386-8d197600ef58', 'WELCOME123', '2026-12-31 23:59:59');

-- ==========================================================
-- 9. TAG RELATIONS
-- ==========================================================
INSERT INTO ticket_announcement_tags (id, tag_id, ticket_id)
VALUES ('10000000-0000-0000-0000-000000000001', 'a0000000-0000-0000-0000-000000000003', 'b0000000-0000-0000-0000-000000000001');

-- +goose Down
DELETE FROM ticket_announcement_tags;
DELETE FROM invite_codes;
DELETE FROM rule_items;
DELETE FROM rules;
DELETE FROM announcements;
DELETE FROM comments;
DELETE FROM tickets;
DELETE FROM users;
DELETE FROM tags;
DELETE FROM apartments WHERE id = '7f23c713-379e-4b9d-8386-8d197600ef58';