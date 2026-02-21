CREATE TABLE IF NOT EXISTS shipments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id VARCHAR(255) NOT NULL,
    tracking_number VARCHAR(255) NOT NULL UNIQUE,
    carrier VARCHAR(100) NOT NULL,
    service VARCHAR(100) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    from_street TEXT NOT NULL,
    from_city VARCHAR(100) NOT NULL,
    from_state VARCHAR(50) NOT NULL,
    from_zip VARCHAR(20) NOT NULL,
    from_country VARCHAR(50) NOT NULL,
    to_street TEXT NOT NULL,
    to_city VARCHAR(100) NOT NULL,
    to_state VARCHAR(50) NOT NULL,
    to_zip VARCHAR(20) NOT NULL,
    to_country VARCHAR(50) NOT NULL,
    weight_grams INT NOT NULL,
    shipping_cost_cents BIGINT NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'USD',
    estimated_days INT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_shipments_order_id ON shipments(order_id);
CREATE INDEX idx_shipments_tracking_number ON shipments(tracking_number);
CREATE INDEX idx_shipments_status ON shipments(status);
