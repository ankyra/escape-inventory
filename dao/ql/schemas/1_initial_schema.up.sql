
CREATE TABLE IF NOT EXISTS release (
    name string, 
    release_id string,
    version string,
    metadata string,
    project string,
    processed_dependencies bool DEFAULT false,
    downloads int DEFAULT 0,
    uploaded_by string DEFAULT "",
    uploaded_at int DEFAULT 0,
);

CREATE UNIQUE INDEX IF NOT EXISTS release_pk ON release (name, version, project);

CREATE TABLE IF NOT EXISTS package (
    project string,
    release_id string, 
    uri string, 
    filesize int DEFAULT -1,
    uploaded_by string DEFAULT "",
    uploaded_at int DEFAULT 0,
);

CREATE UNIQUE INDEX IF NOT EXISTS package_pk ON package (release_id, uri, project);

CREATE TABLE IF NOT EXISTS acl (
    project string,
    group_name string, 
    permission int,
);

CREATE UNIQUE INDEX IF NOT EXISTS acl_pk ON acl (project, group_name);

CREATE TABLE IF NOT EXISTS project (
    name string,
    description string,
    orgURL string,
    logo string,
    hooks string DEFAULT "{}",
);

CREATE UNIQUE INDEX IF NOT EXISTS project_pk ON project (name);

CREATE TABLE IF NOT EXISTS application (
    name string, 
    project string,
    description string DEFAULT "",
    latest_version string DEFAULT "",
    logo string DEFAULT "",
    uploaded_by string DEFAULT "",
    uploaded_at int DEFAULT 0,
    hooks string DEFAULT "{}",
);

CREATE UNIQUE INDEX IF NOT EXISTS application_pk ON application (name, project);

CREATE TABLE IF NOT EXISTS release_dependency (
    project string,
    name string,
    version string,
    dep_project string,
    dep_name string,
    dep_version string,
    build_scope bool,
    deploy_scope bool,
    is_extension bool DEFAULT false,

);

CREATE UNIQUE INDEX IF NOT EXISTS release_dependency_pk ON release_dependency (project, name, version, dep_project, dep_name, dep_version);

CREATE TABLE IF NOT EXISTS subscriptions (
	name string,
	project string,
	subscription_name string,
	subscription_project string,
);

CREATE UNIQUE INDEX subscriptions_pk ON subscriptions (name, project, subscription_name, subscription_project);		 