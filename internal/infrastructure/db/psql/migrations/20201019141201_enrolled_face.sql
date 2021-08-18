-- migrate:up
CREATE SEQUENCE IF NOT EXISTS enrolled_face_sequence
            START WITH 1
            INCREMENT BY 1
            NO MINVALUE
            NO MAXVALUE
            CACHE 1;
CREATE TABLE enrolled_face (
    id BIGINT DEFAULT nextval('enrolled_face_sequence'::regclass) NOT NULL PRIMARY KEY,
    face_id BIGINT NOT NULL,
    name VARCHAR(200),
    deleted_at TIMESTAMP WITH TIME ZONE ,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);
CREATE TRIGGER set_timestamp
    BEFORE UPDATE ON enrolled_face
    FOR EACH ROW
    EXECUTE PROCEDURE trigger_set_timestamp();

-- migrate:down
DROP TRIGGER set_timestamp on enrolled_face;
DROP TABLE enrolled_face;
