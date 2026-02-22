DROP INDEX IF EXISTS idx_addresses_created_at;
ALTER TABLE addresses DROP COLUMN IF EXISTS created_at;
