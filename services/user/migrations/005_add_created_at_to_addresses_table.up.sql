-- Add created_at column for address ordering in existing databases
ALTER TABLE addresses
ADD COLUMN IF NOT EXISTS created_at TIMESTAMP DEFAULT NOW() NOT NULL;

-- Add index for address list ordering
CREATE INDEX IF NOT EXISTS idx_addresses_created_at ON addresses(created_at DESC);
