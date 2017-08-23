CREATE TABLE IF NOT EXISTS release_dependency (
    project VARCHAR(32),
    name VARCHAR(128),
    version VARCHAR(32),
    dep_project VARCHAR(32),
    dep_name VARCHAR(128),
    dep_version VARCHAR(32),
    build_scope BOOLEAN,
    deploy_scope BOOLEAN,
    PRIMARY KEY (project, name, version, 
        dep_project, dep_name, dep_version)
);
