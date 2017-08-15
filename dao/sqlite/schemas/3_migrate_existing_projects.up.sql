INSERT INTO project(name, description, orgURL, logo)
SELECT distinct(project), '', '', ''
FROM release;
