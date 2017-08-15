CREATE TABLE IF NOT EXISTS application (
    name string, 
    project string,
    description string,
    latest_release_id string,
    logo string,
    PRIMARY KEY(name, project)
);
