-- Create wishlists table
CREATE TABLE IF NOT EXISTS wishlists (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    product_id UUID NOT NULL,
    added_at TIMESTAMP DEFAULT NOW() NOT NULL,
    UNIQUE(user_id, product_id)
);

-- Create index on user_id for faster lookups
CREATE INDEX idx_wishlists_user_id ON wishlists(user_id);

-- Create index on product_id
CREATE INDEX idx_wishlists_product_id ON wishlists(product_id);
