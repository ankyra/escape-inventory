CREATE TABLE metrics (
    username string,
    project_count int DEFAULT 0,
);

CREATE UNIQUE INDEX IF NOT EXISTS metrics_pk ON metrics (username);
