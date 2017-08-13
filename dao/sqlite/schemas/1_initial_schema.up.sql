CREATE TABLE IF NOT EXISTS release (
    name string, 
    release_id string,
    version string,
    metadata string,
    project string,
    PRIMARY KEY(name, version, project)
);

CREATE TABLE IF NOT EXISTS package (
    project string,
    release_id string, 
    uri string, 
    PRIMARY KEY(project, release_id, uri)
);

CREATE TABLE IF NOT EXISTS acl (
    project string,
    group_name string, 
    permission varchar(1),
    PRIMARY KEY(project, group_name)
);
