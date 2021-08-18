-- migrate:up
    CREATE SEQUENCE IF NOT EXISTS event
            START WITH 1
            INCREMENT BY 1
            NO MINVALUE
            NO MAXVALUE
            CACHE 1;

    CREATE TABLE event (
        id BIGINT DEFAULT nextval('event'::regclass) NOT NULL,
        type VARCHAR(200) NOT NULL,
        detection JSONB NOT NULL,
        created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
        updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
    ) PARTITION BY RANGE (created_at)
    CREATE INDEX event_auto_index_pk_id on event(id);
	CREATE INDEX event_type_index on event(type);
    CREATE INDEX event_created_at on event(created_at);

    CREATE TRIGGER set_timestamp
        BEFORE UPDATE ON default_credit
        FOR EACH ROW
        EXECUTE PROCEDURE trigger_set_timestamp();
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

-- migrate:down
DROP TRIGGER set_timestamp on enrolled_face;
DROP TABLE event;
