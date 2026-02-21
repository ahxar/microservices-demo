CREATE TABLE IF NOT EXISTS carts (
    user_id TEXT PRIMARY KEY,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS cart_items (
    user_id TEXT NOT NULL,
    product_id TEXT NOT NULL,
    product_name TEXT NOT NULL,
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    unit_price_cents BIGINT NOT NULL CHECK (unit_price_cents >= 0),
    currency TEXT NOT NULL,
    image_url TEXT NOT NULL DEFAULT '',
    PRIMARY KEY (user_id, product_id),
    CONSTRAINT fk_cart_items_cart
        FOREIGN KEY (user_id)
        REFERENCES carts(user_id)
        ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_carts_updated_at ON carts(updated_at);
