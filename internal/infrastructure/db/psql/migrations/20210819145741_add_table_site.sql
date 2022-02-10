-- migrate:up
CREATE SEQUENCE IF NOT EXISTS site_sequence
            START WITH 1
            INCREMENT BY 1
            NO MINVALUE
            NO MAXVALUE
            CACHE 1;
CREATE TABLE site (
    id BIGINT DEFAULT nextval('site_sequence'::regclass) NOT NULL PRIMARY KEY,
    name VARCHAR(200) NOT NUll,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);
CREATE TRIGGER set_timestamp
    BEFORE UPDATE ON site
    FOR EACH ROW
    EXECUTE PROCEDURE trigger_set_timestamp();

-- migrate:down
DROP TABLE site;