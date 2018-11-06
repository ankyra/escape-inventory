CREATE TABLE release_tags (
    project varchar(32),
    application varchar(128), 
    tag varchar(64),
    version varchar(32),
    PRIMARY KEY(project, application, tag)
);
