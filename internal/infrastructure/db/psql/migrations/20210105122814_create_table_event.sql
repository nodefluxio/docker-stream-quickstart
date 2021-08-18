-- migrate:up
CREATE SEQUENCE IF NOT EXISTS event_sequence
        START WITH 1
        INCREMENT BY 1
        NO MINVALUE
        NO MAXVALUE
        CACHE 1;

CREATE TABLE event (
    id BIGINT DEFAULT nextval('event_sequence'::regclass) NOT NULL,
    type VARCHAR(200) NOT NULL,
    stream_id VARCHAR(200) NOT NULL,
    detection JSONB NOT NULL,
    primary_image BYTEA,
    secondary_image BYTEA,
    result JSONB,
    status VARCHAR(200),
    keyword tsvector,
    event_time timestamp with time zone,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
) PARTITION BY RANGE (created_at);

CREATE INDEX event_auto_index_pk_id on event(id);
CREATE INDEX event_type_index on event(type);
CREATE INDEX event_created_at on event(created_at);
CREATE INDEX event_event_time_status on event(event_time, status);

CREATE OR REPLACE FUNCTION create_daily_event(twtz TIMESTAMP WITH TIME ZONE) RETURNS void AS $$
DECLARE
    table_name text := 'event_' || to_char(twtz , 'YYYY_MM_DD');
    start_date text := to_char(twtz , 'YYYY_MM_DD');
    end_date text := to_char(twtz + interval '1 day' , 'YYYY_MM_DD');
BEGIN
    if to_regclass(table_name) IS NULL THEN
        EXECUTE FORMAT('CREATE TABLE %I PARTITION OF event FOR VALUES FROM (%L) to (%L)', table_name, start_date, end_date);
    END IF;
END;
$$ LANGUAGE plpgsql;

SELECT create_daily_event(now()::timestamp);
SELECT create_daily_event(now()::timestamp + interval '1 day');
SELECT create_daily_event(now()::timestamp + interval '2 day');

-- migrate:down
DROP TABLE event;
DROP SEQUENCE IF EXISTS event_sequence;
DROP FUNCTION IF EXISTS create_daily_event;
