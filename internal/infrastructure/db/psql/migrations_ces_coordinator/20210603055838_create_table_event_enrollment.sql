-- migrate:up
CREATE SEQUENCE IF NOT EXISTS event_enrollment_sequence
        START WITH 1
        INCREMENT BY 1
        NO MINVALUE
        NO MAXVALUE
        CACHE 1;

-- partition master table 
CREATE TABLE event_enrollment (
    id BIGINT DEFAULT nextval('event_enrollment_sequence'::regclass) NOT NULL,
    event_id VARCHAR(200) NOT NULL,
    agent VARCHAR(200) NOT NULL,
    event_action VARCHAR(200) NOT NULL,
    payload JSONB NOT NULL,
    event_time timestamp with time zone,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
) PARTITION BY RANGE (created_at);

-- crate index 
CREATE INDEX event_auto_index_pk_id on event_enrollment(id);
CREATE INDEX event_created_at on event_enrollment(created_at);

-- create store procedure function
CREATE OR REPLACE FUNCTION create_daily_event_enrollment(twtz TIMESTAMP WITH TIME ZONE) RETURNS void AS $$
DECLARE
    table_name text := 'event_enrollment_' || to_char(twtz , 'YYYY_MM_DD');
    start_date text := to_char(twtz , 'YYYY_MM_DD');
    end_date text := to_char(twtz + interval '1 day' , 'YYYY_MM_DD');
BEGIN
    if to_regclass(table_name) IS NULL THEN
        EXECUTE FORMAT('CREATE TABLE %I PARTITION OF event_enrollment FOR VALUES FROM (%L) to (%L)', table_name, start_date, end_date);
    END IF;
END;
$$ LANGUAGE plpgsql;

SELECT create_daily_event_enrollment(now()::timestamp);
SELECT create_daily_event_enrollment(now()::timestamp + interval '1 day');
SELECT create_daily_event_enrollment(now()::timestamp + interval '2 day');

-- migrate:down
DROP TABLE event_enrollment;
DROP SEQUENCE IF EXISTS event_enrollment_sequence;
DROP FUNCTION IF EXISTS create_daily_event_enrollment;
