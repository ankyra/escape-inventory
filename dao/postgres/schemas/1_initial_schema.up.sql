CREATE TABLE IF NOT EXISTS release (
    name varchar(128), 
    release_id varchar(256),
    version varchar(32),
    metadata text,
    project varchar(32),
    PRIMARY KEY(name, version, project)
);

CREATE TABLE IF NOT EXISTS package (
    project varchar(32),
    release_id varchar(256), 
    uri varchar(256), 
    PRIMARY KEY(release_id, uri, project)
);

CREATE TABLE IF NOT EXISTS acl (
    project varchar(32),
    group_name varchar(256),
    permission int, 
    PRIMARY KEY(project, group_name)
);
