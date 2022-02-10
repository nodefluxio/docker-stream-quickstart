-- migrate:up
CREATE SEQUENCE IF NOT EXISTS user_sequence
            START WITH 1
            INCREMENT BY 1
            NO MINVALUE
            NO MAXVALUE
            CACHE 1;
CREATE TABLE user_access (
    id BIGINT DEFAULT nextval('user_sequence'::regclass) NOT NULL PRIMARY KEY,
    email VARCHAR(200) NOT NUll,
    username VARCHAR(200) NOT NUll,
    password VARCHAR(200) NOT NUll,
    fullname VARCHAR(200) NOT NUll,
    avatar BYTEA,
    role VARCHAR(200) NOT NUll,
    site_id text[][],
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);
CREATE TRIGGER set_timestamp
    BEFORE UPDATE ON user_access
    FOR EACH ROW
    EXECUTE PROCEDURE trigger_set_timestamp();

INSERT INTO user_access(email, username, password, fullname, role)
VALUES ('admin@admin.com','admin', '$2y$12$RWzFqxS7dCLluJuvxzMXleuiWQd7aehac4SIhwZtYp8.0z3XsbUCOe', 'super admin', 'superadmin');

-- migrate:down
DROP TRIGGER set_timestamp on user_access;
DROP TABLE user_access;
DROP SEQUENCE IF EXISTS user_sequence;
