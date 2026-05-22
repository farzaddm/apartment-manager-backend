-- +goose Up
DROP TABLE IF EXISTS rule_items CASCADE;
DROP TABLE IF EXISTS rules CASCADE;

CREATE TABLE rules (
                       id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                       apartment_id UUID NOT NULL,

                       title VARCHAR(255) NOT NULL,
                       description TEXT DEFAULT '',
                       category rule_category DEFAULT 'other' NOT NULL,

                       created_at TIMESTAMPTZ DEFAULT NOW(),
                       updated_at TIMESTAMPTZ DEFAULT NOW(),
                       deleted_at TIMESTAMPTZ NULL DEFAULT NULL,

                       CONSTRAINT fk_rules_apartment
                           FOREIGN KEY (apartment_id)
                               REFERENCES apartments(id)
                               ON DELETE CASCADE
                               ON UPDATE CASCADE
);

CREATE INDEX idx_rules_apartment_id ON rules(apartment_id);
CREATE INDEX idx_rules_deleted_at ON rules(deleted_at);

CREATE TRIGGER update_rules_updated_at
    BEFORE UPDATE ON rules
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();


-- +goose Down
DROP TABLE IF EXISTS rules CASCADE;

CREATE TABLE rules (
                       id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                       apartment_id UUID NOT NULL,
                       category rule_category DEFAULT 'other' NOT NULL,
                       created_at TIMESTAMPTZ DEFAULT NOW(),
                       updated_at TIMESTAMPTZ DEFAULT NOW(),
                       deleted_at TIMESTAMPTZ NULL DEFAULT NULL,
                       CONSTRAINT fk_rules_apartment FOREIGN KEY (apartment_id) REFERENCES apartments(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE rule_items (
                            id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                            rule_id UUID NOT NULL,
                            body TEXT NOT NULL,
                            created_at TIMESTAMPTZ DEFAULT NOW(),
                            updated_at TIMESTAMPTZ DEFAULT NOW(),
                            deleted_at TIMESTAMPTZ NULL DEFAULT NULL,
                            CONSTRAINT fk_rule_items_rule FOREIGN KEY (rule_id) REFERENCES rules(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE INDEX idx_rules_apartment_id ON rules(apartment_id);
CREATE INDEX idx_rules_deleted_at ON rules(deleted_at);
CREATE TRIGGER update_rules_updated_at BEFORE UPDATE ON rules FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_rule_items_updated_at BEFORE UPDATE ON rule_items FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();