CREATE TABLE release_tags (
    project string,
    application string, 
    tag string,
    version string,
);

CREATE UNIQUE INDEX IF NOT EXISTS release_tags_pk ON release_tags(project, application, tag)
