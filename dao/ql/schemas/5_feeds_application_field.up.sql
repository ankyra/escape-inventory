BEGIN TRANSACTION;
  DROP INDEX feed_events_projects;

  CREATE TABLE tmp_feed_events (
      event_type string,
      username string DEFAULT "",
      project string DEFAULT "",
      application string DEFAULT "",
      timestamp int DEFAULT 0,
      data string DEFAULT "",
  );

  INSERT INTO tmp_feed_events(event_type, username, project, timestamp, data) SELECT event_type, username, project, timestamp, data FROM feed_events;

  DROP TABLE feed_events;
COMMIT;

BEGIN TRANSACTION;
  CREATE TABLE feed_events (
      event_type string,
      username string DEFAULT "",
      project string DEFAULT "",
      application string DEFAULT "",
      timestamp int DEFAULT 0,
      data string DEFAULT "",
  );

  INSERT INTO feed_events(event_type, username, project, application, timestamp, data) SELECT event_type, username, project, application, timestamp, data FROM feed_events;

  DROP TABLE tmp_feed_events;

  CREATE INDEX IF NOT EXISTS feed_events_projects ON feed_events(project);
  CREATE INDEX IF NOT EXISTS feed_events_project_applications ON feed_events(project, application);
COMMIT;
