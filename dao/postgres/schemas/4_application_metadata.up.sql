CREATE TABLE IF NOT EXISTS application (
    name VARCHAR(128), 
    project VARCHAR(32),
    description TEXT DEFAULT '',
    latest_version VARCHAR(256) DEFAULT '',
    logo text DEFAULT '',
    PRIMARY KEY(name, project)
);

INSERT INTO application(project, name)
SELECT DISTINCT (project), name
FROM release;
