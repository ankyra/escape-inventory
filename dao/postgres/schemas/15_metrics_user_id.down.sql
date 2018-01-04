DROP TABLE IF EXISTS metrics;

CREATE TABLE metrics (
    username VARCHAR(64) PRIMARY KEY,
    project_count INT DEFAULT 0
);
