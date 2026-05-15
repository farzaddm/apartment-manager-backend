-- +goose Up


CREATE EXTENSION IF NOT EXISTS "pgcrypto";


-- =========================
-- ENUMS
-- =========================

CREATE TYPE user_role AS ENUM (
    'admin',
    'manager',
    'resident'
);

CREATE TYPE gender_type AS ENUM (
    'male',
    'female'
);

CREATE TYPE ticket_status AS ENUM (
    'open',
    'in-progress',
    'closed'
);

CREATE TYPE ticket_category AS ENUM (
    'maintenance',
    'plumbing',
    'electricity',
    'security',
    'cleaning',
    'parking',
    'other'
);

-- =========================
-- UPDATE updated_at FUNCTION
-- =========================

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $func$
BEGIN
   NEW.updated_at = NOW();
   RETURN NEW;
END;
$func$ LANGUAGE plpgsql;
-- +goose StatementEnd


-- =========================
-- APARTMENTS
-- =========================

CREATE TABLE apartments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    name VARCHAR(255) NOT NULL,

    province VARCHAR(100),
    city VARCHAR(100),
    address TEXT,
    postal_code VARCHAR(20),

    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ NULL DEFAULT NULL
);


-- =========================
-- USERS
-- =========================

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    apartment_id UUID NULL,

    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,

    username VARCHAR(100) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    phone VARCHAR(20) UNIQUE,

    password TEXT NOT NULL,

    role user_role DEFAULT 'resident',

    gender gender_type,

    profile_image_url TEXT,

    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ NULL DEFAULT NULL,

    CONSTRAINT fk_users_apartment
        FOREIGN KEY (apartment_id)
        REFERENCES apartments(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);

-- =========================
-- UNITS
-- =========================

CREATE TABLE units (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    apartment_id UUID NOT NULL,
    user_id UUID NULL,


    unit_number VARCHAR(50) NOT NULL,

    floor INTEGER,

    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ NULL DEFAULT NULL,

    CONSTRAINT fk_units_apartment
        FOREIGN KEY (apartment_id)
        REFERENCES apartments(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE ,

    CONSTRAINT fk_units_user
    FOREIGN KEY (user_id)
    REFERENCES users(id)
    ON DELETE SET NULL 
    ON UPDATE CASCADE
);


-- =========================
-- TICKETS
-- =========================

CREATE TABLE tickets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    user_id UUID  NULL,

    title VARCHAR(255) NOT NULL,
    description TEXT,
    body TEXT,

    category ticket_category DEFAULT 'maintenance',

    status ticket_status DEFAULT 'open',

    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ NULL DEFAULT NULL,

    CONSTRAINT fk_tickets_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE SET NULL 
        ON UPDATE CASCADE
);


-- =========================
-- COMMENTS
-- =========================

CREATE TABLE comments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    user_id UUID  NULL,
    ticket_id UUID NOT NULL,

    body TEXT NOT NULL,
    committed_order INTEGER DEFAULT 0,

    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ NULL DEFAULT NULL,


    CONSTRAINT fk_comments_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE SET NULL 
        ON UPDATE CASCADE,

    CONSTRAINT fk_comments_ticket
        FOREIGN KEY (ticket_id)
        REFERENCES tickets(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);


-- =========================
-- TAGS
-- =========================

CREATE TABLE tags (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    name VARCHAR(100) UNIQUE NOT NULL,

    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ NULL DEFAULT NULL
);


-- =========================
-- ANNOUNCEMENTS
-- =========================

CREATE TABLE announcements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    apartment_id UUID NOT NULL,

    title VARCHAR(255) NOT NULL,
    description TEXT,
    body TEXT,

    expired_date TIMESTAMPTZ NULL,

    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ NULL DEFAULT NULL,

    CONSTRAINT fk_announcements_apartment
        FOREIGN KEY (apartment_id)
        REFERENCES apartments(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);


-- =========================
-- RULES
-- =========================

CREATE TABLE rules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    apartment_id UUID NOT NULL,

    category VARCHAR(255) NOT NULL,

    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ NULL DEFAULT NULL,

    CONSTRAINT fk_rules_apartment
        FOREIGN KEY (apartment_id)
        REFERENCES apartments(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);


-- =========================
-- RULE ITEMS
-- =========================

CREATE TABLE rule_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    rule_id UUID NOT NULL,

    body TEXT NOT NULL,

    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ NULL DEFAULT NULL,

    CONSTRAINT fk_rule_items_rule
        FOREIGN KEY (rule_id)
        REFERENCES rules(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);


-- =========================
-- INVITE CODES
-- =========================

CREATE TABLE invite_codes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    apartment_id UUID NOT NULL,
    unit_id UUID NOT NULL,

    code VARCHAR(64) UNIQUE NOT NULL,

    expires_at TIMESTAMPTZ NOT NULL,

    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ NULL DEFAULT NULL,

    CONSTRAINT fk_invite_codes_apartment
        FOREIGN KEY (apartment_id)
        REFERENCES apartments(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,

    CONSTRAINT fk_invite_codes_unit
        FOREIGN KEY (unit_id)
        REFERENCES units(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);


-- =========================
-- TICKET ANNOUNCEMENT TAGS
-- MANY TO MANY
-- =========================

CREATE TABLE ticket_announcement_tags (
    id UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),

    tag_id UUID NOT NULL,
    ticket_id UUID  NULL,
    announcement_id UUID  NULL,

    
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ NULL DEFAULT NULL,

    CONSTRAINT fk_ticket_announcement_tags_ticket
        FOREIGN KEY (ticket_id)
        REFERENCES tickets(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    
    CONSTRAINT fk_ticket_announcement_tags_announcement
        FOREIGN KEY (announcement_id)
        REFERENCES announcements(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
        
    CONSTRAINT fk_ticket_announcement_tags_tag
        FOREIGN KEY (tag_id)
        REFERENCES tags(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);

-- =========================
-- POLL
-- =========================

CREATE TABLE polls (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    apartment_id UUID NOT NULL,

    title TEXT NOT NULL,
    description TEXT,

    expires_at TIMESTAMPTZ NULL,
    is_votes_public BOOLEAN DEFAULT true ,

    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ NULL DEFAULT NULL,

    CONSTRAINT fk_polls_apartment
    FOREIGN KEY (apartment_id)
    REFERENCES apartments(id)
    ON DELETE CASCADE
    ON UPDATE CASCADE
);

-- =========================
-- poll_OPTIONS
-- =========================

CREATE TABLE poll_options (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    poll_id UUID NOT NULL,

    text TEXT NOT NULL,

    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ NULL DEFAULT NULL,

    CONSTRAINT fk_poll_options_poll
    FOREIGN KEY (poll_id)
    REFERENCES polls(id)
    ON DELETE CASCADE
    ON UPDATE CASCADE
);

-- =========================
-- VOTES
-- =========================

CREATE TABLE votes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    user_id UUID NOT NULL,
    option_id UUID NOT NULL,

    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ NULL DEFAULT NULL,

    CONSTRAINT fk_votes_user
    FOREIGN KEY (user_id)
    REFERENCES users(id)
    ON DELETE CASCADE
    ON UPDATE CASCADE,

    CONSTRAINT fk_votes_option
    FOREIGN KEY (option_id)
    REFERENCES poll_options(id)
    ON DELETE CASCADE
    ON UPDATE CASCADE

);



-- =========================
-- INDEXES
-- =========================

CREATE INDEX idx_users_apartment_id
ON users(apartment_id);

CREATE INDEX idx_users_deleted_at
ON users(deleted_at);

CREATE INDEX idx_tickets_user_id
ON tickets(user_id);

CREATE INDEX idx_tickets_deleted_at
ON tickets(deleted_at);

CREATE INDEX idx_comments_ticket_id
ON comments(ticket_id);

CREATE INDEX idx_comments_user_id
ON comments(user_id);

CREATE INDEX idx_comments_deleted_at
ON comments(deleted_at);

CREATE INDEX idx_announcements_apartment_id
ON announcements(apartment_id);

CREATE INDEX idx_announcements_deleted_at
ON announcements(deleted_at);

CREATE INDEX idx_rules_apartment_id
ON rules(apartment_id);

CREATE INDEX idx_rules_deleted_at
ON rules(deleted_at);

CREATE INDEX idx_invite_codes_apartment_id
ON invite_codes(apartment_id);

CREATE INDEX idx_invite_codes_deleted_at
ON invite_codes(deleted_at);

CREATE INDEX idx_polls_apartment
ON polls(apartment_id);

CREATE INDEX idx_polls_active
ON polls(apartment_id, expires_at)
WHERE deleted_at IS NULL;


CREATE INDEX idx_poll_options_poll
ON poll_options(poll_id)
WHERE deleted_at IS NULL;

CREATE INDEX idx_votes_option
ON votes(option_id)
WHERE deleted_at IS NULL;


CREATE INDEX idx_votes_user
ON votes(user_id)
WHERE deleted_at IS NULL;

CREATE INDEX idx_units_apartment
ON units(apartment_id)
WHERE deleted_at IS NULL;

CREATE INDEX idx_tat_ticket
ON ticket_announcement_tags(ticket_id);

CREATE INDEX idx_tat_announcement
ON ticket_announcement_tags(announcement_id);

CREATE INDEX idx_tat_tag
ON ticket_announcement_tags(tag_id);

CREATE UNIQUE INDEX unique_unit_number_per_apartment
ON units(apartment_id, unit_number)
WHERE deleted_at IS NULL;

-- =========================
-- TRIGGERS
-- =========================

CREATE TRIGGER update_apartments_updated_at
BEFORE UPDATE ON apartments
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_users_updated_at
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_tickets_updated_at
BEFORE UPDATE ON tickets
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_comments_updated_at
BEFORE UPDATE ON comments
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_tags_updated_at
BEFORE UPDATE ON tags
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_announcements_updated_at
BEFORE UPDATE ON announcements
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_rules_updated_at
BEFORE UPDATE ON rules
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_rule_items_updated_at
BEFORE UPDATE ON rule_items
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_invite_codes_updated_at
BEFORE UPDATE ON invite_codes
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_ticket_announcement_tags_updated_at
BEFORE UPDATE ON ticket_announcement_tags
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_unit_updated_at
BEFORE UPDATE ON units
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();


CREATE TRIGGER update_poll_options_updated_at
BEFORE UPDATE ON poll_options
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_polls_updated_at
BEFORE UPDATE ON polls
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_votes_updated_at
BEFORE UPDATE ON votes
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();



-- +goose Down
DROP TABLE IF EXISTS invite_codes CASCADE;

DROP TABLE IF EXISTS rule_items CASCADE;
DROP TABLE IF EXISTS rules CASCADE;

DROP TABLE IF EXISTS announcements CASCADE;

DROP TABLE IF EXISTS tags CASCADE;

DROP TABLE IF EXISTS comments CASCADE;

DROP TABLE IF EXISTS tickets CASCADE;

DROP TABLE IF EXISTS users CASCADE;

DROP TABLE IF EXISTS apartments CASCADE;

DROP TABLE IF EXISTS ticket_announcement_tags CASCADE;

DROP TABLE IF EXISTS units CASCADE;

DROP TABLE IF EXISTS votes CASCADE;

DROP TABLE IF EXISTS poll_options CASCADE;

DROP TABLE IF EXISTS polls CASCADE;


-- =========================================
-- DROP ENUMS
-- =========================================

DROP TYPE IF EXISTS ticket_status CASCADE;
DROP TYPE IF EXISTS gender_type CASCADE;
DROP TYPE IF EXISTS user_role CASCADE;
DROP TYPE IF EXISTS ticket_category CASCADE;



-- =========================================
-- DROP TRIGGERS
-- =========================================

DROP TRIGGER IF EXISTS update_apartments_updated_at ON apartments;
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
DROP TRIGGER IF EXISTS update_tickets_updated_at ON tickets;
DROP TRIGGER IF EXISTS update_comments_updated_at ON comments;
DROP TRIGGER IF EXISTS update_tags_updated_at ON tags;
DROP TRIGGER IF EXISTS update_announcements_updated_at ON announcements;
DROP TRIGGER IF EXISTS update_rules_updated_at ON rules;
DROP TRIGGER IF EXISTS update_rule_items_updated_at ON rule_items;
DROP TRIGGER IF EXISTS update_invite_codes_updated_at ON invite_codes;
DROP TRIGGER IF EXISTS update_ticket_announcement_tags_updated_at ON ticket_announcement_tags;
DROP TRIGGER IF EXISTS update_unit_updated_at ON units;
DROP TRIGGER IF EXISTS update_votes_updated_at ON votes;
DROP TRIGGER IF EXISTS update_poll_options_updated_at ON poll_options;
DROP TRIGGER IF EXISTS update_polls_updated_at ON polls;





-- =========================================
-- DROP FUNCTIONS
-- =========================================

DROP FUNCTION IF EXISTS update_updated_at_column CASCADE;


-- =========================================
-- DROP EXTENSION
-- =========================================

DROP EXTENSION IF EXISTS "pgcrypto";