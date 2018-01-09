CREATE TABLE feed_events (
    id SERIAL PRIMARY KEY,
    event_type VARCHAR(32),
    username VARCHAR(64),
    project VARCHAR(32),
    timestamp INTEGER DEFAULT '0',
    data BYTEA
);
