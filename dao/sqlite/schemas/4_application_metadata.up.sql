CREATE TABLE IF NOT EXISTS application (
    name string, 
    project string,
    description string default '',
    latest_version string default '',
    logo string default '',
    PRIMARY KEY(name, project)
);

INSERT INTO application(project, name)
SELECT DISTINCT (project), name
FROM release;
