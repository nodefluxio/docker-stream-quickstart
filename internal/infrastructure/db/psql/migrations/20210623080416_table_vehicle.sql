-- migrate:up
CREATE SEQUENCE IF NOT EXISTS vehicle_sequence
            START WITH 1
            INCREMENT BY 1
            NO MINVALUE
            NO MAXVALUE
            CACHE 1;
CREATE TABLE vehicle (
    id BIGINT DEFAULT nextval('vehicle_sequence'::regclass) NOT NULL PRIMARY KEY,
    plate_number VARCHAR(100) NOT NULL,
    unique_id varchar(100) UNIQUE NOT NULL,
    type varchar(100),
    brand varchar(100),
    color varchar(100),
    name VARCHAR(200),
    status varchar(100),
    deleted_at TIMESTAMP WITH TIME ZONE ,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX vehicle_auto_index_id on vehicle(id);
CREATE INDEX vehicle_index_name on vehicle(name);
CREATE INDEX vehicle_index_plate_number on vehicle(plate_number);
CREATE INDEX vehicle_index_status on vehicle(status);
CREATE INDEX vehicle_index_type on vehicle(type);
CREATE INDEX vehicle_index_unique_id on vehicle(unique_id);

CREATE TRIGGER set_timestamp
    BEFORE UPDATE ON vehicle
    FOR EACH ROW
    EXECUTE PROCEDURE trigger_set_timestamp();

-- migrate:down
DROP TRIGGER set_timestamp on vehicle;
DROP TABLE vehicle;
