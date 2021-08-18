-- migrate:up
CREATE SEQUENCE IF NOT EXISTS global_setting_sequence
            START WITH 1
            INCREMENT BY 1
            NO MINVALUE
            NO MAXVALUE
            CACHE 1;
CREATE TABLE global_setting (
    id BIGINT DEFAULT nextval('global_setting_sequence'::regclass) NOT NULL PRIMARY KEY,
    similarity DECIMAL NOT NUll,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- migrate:down
DROP TABLE global_setting;
