CREATE TABLE providers (
    project string,
    application string, 
    version string, 
    description string,
    provider string
);

CREATE UNIQUE INDEX IF NOT EXISTS provider_pk ON providers(project, application, provider);
