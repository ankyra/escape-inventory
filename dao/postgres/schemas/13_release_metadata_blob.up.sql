ALTER TABLE release DROP COLUMN metadata;
ALTER TABLE release ADD COLUMN metadata bytea;
