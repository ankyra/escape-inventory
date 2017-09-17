CREATE TABLE subscriptions (
  name string,
  project string,
  subscription_name string,
  subscription_project string,
  PRIMARY KEY (name, project, subscription_name, subscription_project)
);
