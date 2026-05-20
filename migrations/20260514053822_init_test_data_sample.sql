-- +goose Up
-- ============================================================================
-- SEED DATA
-- ============================================================================

-- 1. Insert 1 Apartment
INSERT INTO apartments (id, name, province, city, address, postal_code)
VALUES (
    'a0000000-0000-0000-0000-000000000000', 
    'Grand Horizon Luxury Apartments', 
    'Ontario', 
    'Toronto', 
    '123 Skyview Avenue, Suite 100', 
    'M5V 2T6'
);

-- 2. Insert Users (1 Admin, 1 Manager, 4 Residents)
INSERT INTO users (id, apartment_id, first_name, last_name, username, email, phone, password, role, gender)
VALUES 
(
    '11111111-1111-1111-1111-111111111111', 
    'a0000000-0000-0000-0000-000000000000', 
    'Alice', 'Smith', 'alice_admin', 'alice@horizon.com', '+14165550001', 
    '$2b$12$LRY8b1Bf2h...', 'admin', 'female'
),
(
    '22222222-2222-2222-2222-222222222222', 
    'a0000000-0000-0000-0000-000000000000', 
    'Bob', 'Jones', 'bob_manager', 'bob@horizon.com', '+14165550002', 
    '$2b$12$LRY8b1Bf2h...', 'manager', 'male'
),
(
    '33333333-3333-3333-3333-333333333333', 
    'a0000000-0000-0000-0000-000000000000', 
    'Charlie', 'Brown', 'charlie_r1', 'charlie@gmail.com', '+14165550003', 
    '$2b$12$LRY8b1Bf2h...', 'resident', 'male'
),
(
    '44444444-4444-4444-4444-444444444444', 
    'a0000000-0000-0000-0000-000000000000', 
    'Diana', 'Prince', 'diana_r2', 'diana@gmail.com', '+14165550004', 
    '$2b$12$LRY8b1Bf2h...', 'resident', 'female'
),
(
    '55555555-5555-5555-5555-555555555555', 
    'a0000000-0000-0000-0000-000000000000', 
    'Evan', 'Wright', 'evan_r3', 'evan@gmail.com', '+14165550005', 
    '$2b$12$LRY8b1Bf2h...', 'resident', 'male'
),
(
    '66666666-6666-6666-6666-666666666666', 
    'a0000000-0000-0000-0000-000000000000', 
    'Fiona', 'Gallagher', 'fiona_r4', 'fiona@gmail.com', '+14165550006', 
    '$2b$12$LRY8b1Bf2h...', 'resident', 'female'
),
(
    '77777777-7777-7777-7777-777777777777', 
    'a0000000-0000-0000-0000-000000000000', 
    'Voldemort', 'VoldemortianResident', 'voldemort_res', 'voldemort@gmail.com', '+14165550097', 
    '$2b$12$LRY8b1Bf2h...', 'resident', 'male'
),
(
    '88888888-8888-8888-8888-888888888888', 
    'a0000000-0000-0000-0000-000000000000', 
    'Voldemort', 'VoldemortianManager', 'voldemort_man', 'voldemort@gmail.com', '+14165550999', 
    '$2b$12$LRY8b1Bf2h...', 'manager', 'male'
),
(
    '99999999-9999-9999-9999-999999999999', 
    'a0000000-0000-0000-0000-000000000000', 
    'Voldemort', 'VoldemortianAdmin', 'voldemort_ad', 'voldemort@gmail.com', '+14165559797', 
    '$2b$12$LRY8b1Bf2h...', 'admin', 'male'
);

-- 3. Insert Units (6 units assigned to various users/empty)
INSERT INTO units (id, apartment_id, user_id, unit_number, floor)
VALUES 
('11111111-eeee-eeee-eeee-eeeeeeeeeeee', 'a0000000-0000-0000-0000-000000000000', '33333333-3333-3333-3333-333333333333', '101', 1),
('22222222-eeee-eeee-eeee-eeeeeeeeeeee', 'a0000000-0000-0000-0000-000000000000', '44444444-4444-4444-4444-444444444444', '102', 1),
('33333333-eeee-eeee-eeee-eeeeeeeeeeee', 'a0000000-0000-0000-0000-000000000000', '55555555-5555-5555-5555-555555555555', '201', 2),
('44444444-eeee-eeee-eeee-eeeeeeeeeeee', 'a0000000-0000-0000-0000-000000000000', '66666666-6666-6666-6666-666666666666', '202', 2),
('55555555-eeee-eeee-eeee-eeeeeeeeeeee', 'a0000000-0000-0000-0000-000000000000', NULL, '301', 3),
('66666666-eeee-eeee-eeee-eeeeeeeeeeee', 'a0000000-0000-0000-0000-000000000000', NULL, '302', 3);

-- 4. Insert Tickets (6 maintenance/complaint tickets)
INSERT INTO tickets (id, user_id, title, description, body, category, status)
VALUES 
(
    '11111111-bbbb-bbbb-bbbb-bbbbbbbbbbbb', '33333333-3333-3333-3333-333333333333',
    'Leaking Kitchen Sink', 'Water dripping under sink', 'The pipe underneath the kitchen sink has a slow drip.', 'plumbing', 'open'
),
(
    '22222222-bbbb-bbbb-bbbb-bbbbbbbbbbbb', '44444444-4444-4444-4444-444444444444',
    'AC Not Cooling', 'AC blowing warm air', 'AC fan works but air is not cold. High temperature in room.', 'maintenance', 'in-progress'
),
(
    '33333333-bbbb-bbbb-bbbb-bbbbbbbbbbbb', '55555555-5555-5555-5555-555555555555',
    'Broken Front Door Lock', 'Electronic keypad failing', 'Keypad takes 4-5 tries to recognize code.', 'security', 'open'
),
(
    '44444444-bbbb-bbbb-bbbb-bbbbbbbbbbbb', '66666666-6666-6666-6666-666666666666',
    'Loud Noise Complaint', 'Neighbor playing bass late at night', 'Unit 201 has been playing heavy electronic music past 11 PM.', 'other', 'closed'
),
(
    '55555555-bbbb-bbbb-bbbb-bbbbbbbbbbbb', '33333333-3333-3333-3333-333333333333',
    'Gym Equipment Broken', 'Treadmill #2 belt loose', 'The treadmill closest to the window slips when running.', 'maintenance', 'open'
),
(
    '66666666-bbbb-bbbb-bbbb-bbbbbbbbbbbb', '44444444-4444-4444-4444-444444444444',
    'Light Bulb Replacement', 'Hallway light flickers', 'The light outside my door is flickering constantly.', 'electricity', 'closed'
);

-- 5. Insert Comments (6 comments on the tickets)
INSERT INTO comments (id, user_id, ticket_id, body, committed_order)
VALUES 
('11111111-cccc-cccc-cccc-cccccccccccc', '22222222-2222-2222-2222-222222222222', '22222222-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 'HVAC technician has been scheduled for tomorrow morning.', 1),
('22222222-cccc-cccc-cccc-cccccccccccc', '44444444-4444-4444-4444-444444444444', '22222222-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 'Thanks for the quick response!', 2),
('33333333-cccc-cccc-cccc-cccccccccccc', '22222222-2222-2222-2222-222222222222', '44444444-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 'Spoke to the resident. Issue should be resolved.', 1),
('44444444-cccc-cccc-cccc-cccccccccccc', '66666666-6666-6666-6666-666666666666', '44444444-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 'Confirmed, it is quiet now. Thank you.', 2),
('55555555-cccc-cccc-cccc-cccccccccccc', '22222222-2222-2222-2222-222222222222', '66666666-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 'Bulb replaced by building super.', 1),
('66666666-cccc-cccc-cccc-cccccccccccc', '22222222-2222-2222-2222-222222222222', '11111111-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 'Plumber is looking at this on Thursday.', 1);

-- 6. Insert Tags (6 tags)
INSERT INTO tags (id, name)
VALUES 
('11111111-dddd-dddd-dddd-dddddddddddd', 'Urgent'),
('22222222-dddd-dddd-dddd-dddddddddddd', 'Maintenance'),
('33333333-dddd-dddd-dddd-dddddddddddd', 'Notice'),
('44444444-dddd-dddd-dddd-dddddddddddd', 'Event'),
('55555555-dddd-dddd-dddd-dddddddddddd', 'Plumbing'),
('66666666-dddd-dddd-dddd-dddddddddddd', 'Safety');

-- 7. Insert Announcements (6 announcements)
INSERT INTO announcements (id, apartment_id, title, description, body, expired_date)
VALUES 
(
    '11111111-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'a0000000-0000-0000-0000-000000000000',
    'Water Shutoff Notice', 'Scheduled maintenance work', 'Water will be turned off from 9 AM to 1 PM this Friday.', NOW() + INTERVAL '2 days'
),
(
    '22222222-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'a0000000-0000-0000-0000-000000000000',
    'Roof Inspection', 'Contractors on site', 'Please expect workers on the roof throughout the upcoming week.', NOW() + INTERVAL '7 days'
),
(
    '33333333-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'a0000000-0000-0000-0000-000000000000',
    'Summer BBQ Party', 'Community event', 'Join us on the rooftop terrace this Saturday at 5 PM for free food and drinks.', NOW() + INTERVAL '3 days'
),
(
    '44444444-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'a0000000-0000-0000-0000-000000000000',
    'Fire Alarm Testing', 'Routine monthly check', 'Alarms will sound intermittently between 10 AM and 12 PM.', NOW() + INTERVAL '1 day'
),
(
    '55555555-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'a0000000-0000-0000-0000-000000000000',
    'Elevator Modernization Update', 'West Elevator offline', 'West elevator upgrades will take an extra 3 days to complete.', NOW() + INTERVAL '4 days'
),
(
    '66666666-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'a0000000-0000-0000-0000-000000000000',
    'Parking Lot Cleaning', 'Move your vehicles', 'Lower deck parking will be power-washed. All vehicles must be moved.', NOW() + INTERVAL '5 days'
);

-- 8. Insert Rules (6 rule categories)
INSERT INTO rules (id, apartment_id, category)
VALUES 
('11111111-ffff-ffff-ffff-ffffffffffff', 'a0000000-0000-0000-0000-000000000000', 'pet_policy'),
('22222222-ffff-ffff-ffff-ffffffffffff', 'a0000000-0000-0000-0000-000000000000', 'noise_regulations'),
('33333333-ffff-ffff-ffff-ffffffffffff', 'a0000000-0000-0000-0000-000000000000', 'gym_rules'),
('44444444-ffff-ffff-ffff-ffffffffffff', 'a0000000-0000-0000-0000-000000000000', 'garbage_recycling'),
('55555555-ffff-ffff-ffff-ffffffffffff', 'a0000000-0000-0000-0000-000000000000', 'parking_bylaws'),
('66666666-ffff-ffff-ffff-ffffffffffff', 'a0000000-0000-0000-0000-000000000000', 'pool_policy');

-- 9. Insert Rule Items (6 details mapped to rules)
INSERT INTO rule_items (id, rule_id, body)
VALUES 
('11111111-2222-3333-4444-555555555551', '11111111-ffff-ffff-ffff-ffffffffffff', 'Pets must be kept on a leash at all times in common areas.'),
('11111111-2222-3333-4444-555555555552', '22222222-ffff-ffff-ffff-ffffffffffff', 'Quiet hours are from 11:00 PM to 7:00 AM daily.'),
('11111111-2222-3333-4444-555555555553', '33333333-ffff-ffff-ffff-ffffffffffff', 'Please wipe down gym machines after usage.'),
('11111111-2222-3333-4444-555555555554', '44444444-ffff-ffff-ffff-ffffffffffff', 'Cardboard boxes must be flattened before being placed in chutes.'),
('11111111-2222-3333-4444-555555555555', '55555555-ffff-ffff-ffff-ffffffffffff', 'Visitor parking is restricted to a maximum of 24 consecutive hours.'),
('11111111-2222-3333-4444-555555555556', '66666666-ffff-ffff-ffff-ffffffffffff', 'No glass containers allowed in the pool enclosure.');

-- 10. Insert Invite Codes (6 codes for pending residents/units)
INSERT INTO invite_codes (id, apartment_id, unit_id, code, expires_at)
VALUES 
('11111111-9999-9999-9999-999999999991', 'a0000000-0000-0000-0000-000000000000', '55555555-eeee-eeee-eeee-eeeeeeeeeeee', 'INV-301-WELCOME', NOW() + INTERVAL '30 days'),
('11111111-9999-9999-9999-999999999992', 'a0000000-0000-0000-0000-000000000000', '66666666-eeee-eeee-eeee-eeeeeeeeeeee', 'INV-302-WELCOME', NOW() + INTERVAL '30 days'),
('11111111-9999-9999-9999-999999999993', 'a0000000-0000-0000-0000-000000000000', '11111111-eeee-eeee-eeee-eeeeeeeeeeee', 'INV-EXPIRED-99', NOW() - INTERVAL '1 day'),
('11111111-9999-9999-9999-999999999994', 'a0000000-0000-0000-0000-000000000000', '22222222-eeee-eeee-eeee-eeeeeeeeeeee', 'CODE-ALPHA-444', NOW() + INTERVAL '10 days'),
('11111111-9999-9999-9999-999999999995', 'a0000000-0000-0000-0000-000000000000', '33333333-eeee-eeee-eeee-eeeeeeeeeeee', 'CODE-BETA-555',  NOW() + INTERVAL '15 days'),
('11111111-9999-9999-9999-999999999996', 'a0000000-0000-0000-0000-000000000000', '44444444-eeee-eeee-eeee-eeeeeeeeeeee', 'CODE-GAMMA-666', NOW() + INTERVAL '20 days');

-- 11. Insert Many-To-Many Relations (ticket_announcement_tags)
INSERT INTO ticket_announcement_tags (id, tag_id, ticket_id, announcement_id)
VALUES 
('11111111-8888-8888-8888-888888888881', '11111111-dddd-dddd-dddd-dddddddddddd', '11111111-bbbb-bbbb-bbbb-bbbbbbbbbbbb', NULL), 
('11111111-8888-8888-8888-888888888882', '55555555-dddd-dddd-dddd-dddddddddddd', '11111111-bbbb-bbbb-bbbb-bbbbbbbbbbbb', NULL), 
('11111111-8888-8888-8888-888888888883', '11111111-dddd-dddd-dddd-dddddddddddd', NULL, '11111111-aaaa-aaaa-aaaa-aaaaaaaaaaaa'), 
('11111111-8888-8888-8888-888888888884', '33333333-dddd-dddd-dddd-dddddddddddd', NULL, '11111111-aaaa-aaaa-aaaa-aaaaaaaaaaaa'), 
('11111111-8888-8888-8888-888888888885', '44444444-dddd-dddd-dddd-dddddddddddd', NULL, '33333333-aaaa-aaaa-aaaa-aaaaaaaaaaaa'), 
('11111111-8888-8888-8888-888888888886', '66666666-dddd-dddd-dddd-dddddddddddd', '33333333-bbbb-bbbb-bbbb-bbbbbbbbbbbb', NULL);


-- 12. Insert Polls (2 distinct polls for the apartment)
INSERT INTO polls (id, apartment_id, title, description, expires_at, is_votes_public)
VALUES 
(
    '11111111-7777-7777-7777-777777777771', 
    'a0000000-0000-0000-0000-000000000000', 
    'Rooftop BBQ Painting Color', 
    'Please vote on the new accent wall color scheme for the rooftop terrace.', 
    NOW() + INTERVAL '14 days', 
    true
),
(
    '11111111-7777-7777-7777-777777777772', 
    'a0000000-0000-0000-0000-000000000000', 
    'Gym Hour Extension', 
    'Should we extend gym operating hours to 24/7? (Anonymous tracking verification)', 
    NOW() + INTERVAL '7 days', 
    false
);

-- 13. Insert Options (3 options for the BBQ poll, 2 options for the Gym poll)
INSERT INTO poll_options (id, poll_id, text)
VALUES 
-- Options for Poll 1 (BBQ Paint)
('11111111-3333-4444-5555-6666666666a1', '11111111-7777-7777-7777-777777777771', 'Ocean Breeze Blue'),
('11111111-3333-4444-5555-6666666666a2', '11111111-7777-7777-7777-777777777771', 'Urban Charcoal Gray'),
('11111111-3333-4444-5555-6666666666a3', '11111111-7777-7777-7777-777777777771', 'Terracotta Sunset'),

-- Options for Poll 2 (Gym Hours)
('11111111-3333-4444-5555-6666666666b1', '11111111-7777-7777-7777-777777777772', 'Yes, extend to 24/7'),
('11111111-3333-4444-5555-6666666666b2', '11111111-7777-7777-7777-777777777772', 'No, keep current hours (6 AM - 11 PM)');

-- 14. Insert Votes (6 structural user votes mapped to choices)
INSERT INTO votes (id, user_id, option_id)
VALUES 
-- Votes for Poll 1 (BBQ Paint)
('11111111-0000-1111-2222-333333333331', '33333333-3333-3333-3333-333333333333', '11111111-3333-4444-5555-6666666666a1'), -- Charlie votes Blue
('11111111-0000-1111-2222-333333333332', '44444444-4444-4444-4444-444444444444', '11111111-3333-4444-5555-6666666666a2'), -- Diana votes Gray
('11111111-0000-1111-2222-333333333333', '55555555-5555-5555-5555-555555555555', '11111111-3333-4444-5555-6666666666a1'), -- Evan votes Blue

-- Votes for Poll 2 (Gym Hours)
('11111111-0000-1111-2222-333333333334', '33333333-3333-3333-3333-333333333333', '11111111-3333-4444-5555-6666666666b1'), -- Charlie votes Yes
('11111111-0000-1111-2222-333333333335', '44444444-4444-4444-4444-444444444444', '11111111-3333-4444-5555-6666666666b1'), -- Diana votes Yes
('11111111-0000-1111-2222-333333333336', '66666666-6666-6666-6666-666666666666', '11111111-3333-4444-5555-6666666666b2'); -- Fiona votes No



-- +goose Down
-- ============================================================================
-- TRUNCATE SEED DATA
-- ============================================================================
TRUNCATE TABLE votes CASCADE;
TRUNCATE TABLE poll_options CASCADE;
TRUNCATE TABLE polls CASCADE;
TRUNCATE TABLE ticket_announcement_tags CASCADE;
TRUNCATE TABLE invite_codes CASCADE;
TRUNCATE TABLE rule_items CASCADE;
TRUNCATE TABLE rules CASCADE;
TRUNCATE TABLE announcements CASCADE;
TRUNCATE TABLE tags CASCADE;
TRUNCATE TABLE comments CASCADE;
TRUNCATE TABLE tickets CASCADE;
TRUNCATE TABLE units CASCADE;
TRUNCATE TABLE users CASCADE;
TRUNCATE TABLE apartments CASCADE;