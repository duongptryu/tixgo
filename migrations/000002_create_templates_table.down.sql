-- Drop indexes
DROP INDEX IF EXISTS idx_templates_created_at;
DROP INDEX IF EXISTS idx_templates_created_by;
DROP INDEX IF EXISTS idx_templates_status;
DROP INDEX IF EXISTS idx_templates_type;
DROP INDEX IF EXISTS idx_templates_slug;

-- Drop foreign key constraint if it was created
-- ALTER TABLE templates DROP CONSTRAINT IF EXISTS fk_templates_created_by;

-- Drop templates table
DROP TABLE IF EXISTS templates; 