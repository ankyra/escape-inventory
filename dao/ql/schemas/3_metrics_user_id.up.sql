DROP INDEX metrics_pk;
DROP TABLE metrics;

CREATE TABLE metrics (
    user_id string,
    project_count int DEFAULT 0,
);

CREATE UNIQUE INDEX IF NOT EXISTS metrics_pk ON metrics (user_id);
