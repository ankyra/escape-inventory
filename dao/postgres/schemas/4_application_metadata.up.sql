CREATE TABLE IF NOT EXISTS application (
    name VARCHAR(128), 
    project VARCHAR(32),
    description TEXT DEFAULT '',
    latest_release_id VARCHAR(256) DEFAULT '',
    logo text DEFAULT '',
    PRIMARY KEY(name, project)
);
