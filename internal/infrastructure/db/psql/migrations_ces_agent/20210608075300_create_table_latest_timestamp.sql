-- migrate:up
CREATE TABLE latest_timestamp (
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- migrate:down
DROP TABLE latest_timestamp;
