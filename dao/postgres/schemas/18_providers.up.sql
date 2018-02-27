CREATE TABLE providers (
    project varchar(32),
    application varchar(128), 
    version varchar(32),
    description TEXT DEFAULT '',
    provider varchar(32),
    PRIMARY KEY(project, application, provider)
);

