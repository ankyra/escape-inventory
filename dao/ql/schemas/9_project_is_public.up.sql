BEGIN TRANSACTION;
	ALTER TABLE project ADD is_public bool DEFAULT false;
COMMIT;

BEGIN TRANSACTION;
  	DROP INDEX project_pk;

	CREATE TABLE IF NOT EXISTS tmp_project (
		name string,
		description string,
		orgURL string,
		logo string,
		hooks string DEFAULT "{}",
		is_public bool DEFAULT false,
	);

  	INSERT INTO tmp_project(name, description, orgURL, logo, hooks) SELECT name, description, orgURL, logo, hooks FROM project;

 	DROP TABLE project;
COMMIT;

BEGIN TRANSACTION;
	CREATE TABLE IF NOT EXISTS project (
		name string,
		description string,
		orgURL string,
		logo string,
		hooks string DEFAULT "{}",
		is_public bool DEFAULT false,
	);

  	INSERT INTO project(name, description, orgURL, logo, hooks) SELECT name, description, orgURL, logo, hooks FROM tmp_project;

  	DROP TABLE tmp_project;

	CREATE UNIQUE INDEX IF NOT EXISTS project_pk ON project (name);
COMMIT;