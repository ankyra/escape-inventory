CREATE TABLE subscriptions(
    project VARCHAR(32),
    name VARCHAR(128),
    subscription_project VARCHAR(32),
    subscription_name VARCHAR(128),
    PRIMARY KEY(project, name, subscription_project, subscription_name)
);
