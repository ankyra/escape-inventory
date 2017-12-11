ALTER TABLE release ALTER COLUMN metadata TYPE bytea USING metadata::bytea;
