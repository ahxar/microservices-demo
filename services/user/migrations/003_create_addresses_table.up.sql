-- Create addresses table
CREATE TABLE IF NOT EXISTS addresses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    label VARCHAR(50),
    street VARCHAR(255) NOT NULL,
    city VARCHAR(100) NOT NULL,
    state VARCHAR(100),
    zip_code VARCHAR(20) NOT NULL,
    country VARCHAR(100) NOT NULL,
    is_default BOOLEAN DEFAULT false NOT NULL
);

-- Create index on user_id for faster lookups
CREATE INDEX idx_addresses_user_id ON addresses(user_id);

-- Create index on is_default for filtering
CREATE INDEX idx_addresses_is_default ON addresses(user_id, is_default);
