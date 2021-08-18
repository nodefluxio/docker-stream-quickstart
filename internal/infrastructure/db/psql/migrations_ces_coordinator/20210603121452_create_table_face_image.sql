-- migrate:up
CREATE SEQUENCE IF NOT EXISTS face_image_sequence
            START WITH 1
            INCREMENT BY 1
            NO MINVALUE
            NO MAXVALUE
            CACHE 1;

CREATE TABLE face_image (
    id BIGINT DEFAULT nextval('face_image_sequence'::regclass) NOT NULL PRIMARY KEY,
    event_enrollment_id BIGINT NOT NULL,
    image BYTEA NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- migrate:down
DROP TABLE face_image;
