-- migrate:up
CREATE SEQUENCE IF NOT EXISTS map_site_stream_sequence
            START WITH 1
            INCREMENT BY 1
            NO MINVALUE
            NO MAXVALUE
            CACHE 1;
CREATE TABLE map_site_stream (
    id BIGINT DEFAULT nextval('map_site_stream_sequence'::regclass) NOT NULL PRIMARY KEY,
    site_id BIGINT REFERENCES site(id) NOT NUll,
    stream_id VARCHAR(200) NOT NUll,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- migrate:down
DROP TABLE map_site_stream;