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

CREATE TABLE schema_migrations (
    version bigint NOT NULL,
    dirty boolean NOT NULL
);

ALTER TABLE schema_migrations OWNER TO postgres;

INSERT INTO schema_migrations VALUES (1, false);

INSERT INTO release(name, release_id, version, metadata, project) VALUES (
    'test', 
    'project/test-v1.0', 
    '1.0', 
    '{"name": "test", "version": "1.0", "project": "project", "description": "yo"}', 
    'project'
);
