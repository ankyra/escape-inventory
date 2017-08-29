DELETE FROM release_dependency;
UPDATE release SET processed_dependencies = 'false';
ALTER TABLE release_dependency ADD COLUMN is_extension BOOLEAN DEFAULT 'false';
