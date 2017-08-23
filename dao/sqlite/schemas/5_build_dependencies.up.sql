CREATE TABLE IF NOT EXISTS release_dependency (
    project string,
    name string,
    version string,
    dep_project string,
    dep_name string,
    dep_version string,
    build_scope boolean,
    deploy_scope boolean,
    PRIMARY KEY (project, name, version, dep_project, dep_name, dep_version)
);
