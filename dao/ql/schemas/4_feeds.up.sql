CREATE TABLE feed_events (
    event_type string,
    username string DEFAULT "",
    project string DEFAULT "",
    timestamp int DEFAULT 0,
    data string DEFAULT "",
);

CREATE INDEX IF NOT EXISTS feed_events_projects ON feed_events(project);
