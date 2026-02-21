#!/bin/bash

# Seed sample catalog data directly into PostgreSQL.
# This script is idempotent and can be re-run safely.

set -euo pipefail

POSTGRES_CONTAINER="microservices-postgres"
POSTGRES_USER="postgres"
CATALOG_DB="catalog_db"

echo "🌱 Seeding sample catalog data..."
echo ""

if ! docker ps --format '{{.Names}}' | grep -q "^${POSTGRES_CONTAINER}$"; then
  echo "❌ Error: PostgreSQL container '${POSTGRES_CONTAINER}' is not running."
  echo "Start services first:"
  echo "  docker compose -f deployments/docker/docker-compose.yml up -d postgres"
  exit 1
fi

docker exec -i "${POSTGRES_CONTAINER}" psql -U "${POSTGRES_USER}" -d "${CATALOG_DB}" <<'SQL'
-- Categories
INSERT INTO categories (name, slug, description)
VALUES
  ('Electronics', 'electronics', 'Electronic devices and gadgets'),
  ('Clothing', 'clothing', 'Fashion and apparel'),
  ('Books', 'books', 'Books and literature')
ON CONFLICT (slug) DO UPDATE
SET
  name = EXCLUDED.name,
  description = EXCLUDED.description;

-- Products
INSERT INTO products (name, slug, description, price_cents, currency, category_id, image_urls, stock_quantity, is_active)
VALUES
  ('Wireless Headphones', 'wireless-headphones', 'Premium noise-canceling wireless headphones with 30-hour battery life', 29900, 'USD', (SELECT id FROM categories WHERE slug = 'electronics'), ARRAY[]::TEXT[], 50, true),
  ('Smart Watch', 'smart-watch', 'Fitness tracking smartwatch with heart rate monitor and GPS', 39900, 'USD', (SELECT id FROM categories WHERE slug = 'electronics'), ARRAY[]::TEXT[], 30, true),
  ('Laptop Backpack', 'laptop-backpack', 'Durable backpack with padded laptop compartment, fits up to 15 inch laptops', 5900, 'USD', (SELECT id FROM categories WHERE slug = 'electronics'), ARRAY[]::TEXT[], 100, true),
  ('Cotton T-Shirt', 'cotton-t-shirt', 'Comfortable 100% cotton t-shirt, available in multiple colors', 1999, 'USD', (SELECT id FROM categories WHERE slug = 'clothing'), ARRAY[]::TEXT[], 200, true),
  ('Denim Jeans', 'denim-jeans', 'Classic fit denim jeans, durable and stylish', 4999, 'USD', (SELECT id FROM categories WHERE slug = 'clothing'), ARRAY[]::TEXT[], 75, true),
  ('Winter Jacket', 'winter-jacket', 'Warm and waterproof winter jacket with hood', 12900, 'USD', (SELECT id FROM categories WHERE slug = 'clothing'), ARRAY[]::TEXT[], 40, true),
  ('The Great Gatsby', 'the-great-gatsby', 'Classic American novel by F. Scott Fitzgerald', 1499, 'USD', (SELECT id FROM categories WHERE slug = 'books'), ARRAY[]::TEXT[], 150, true),
  ('To Kill a Mockingbird', 'to-kill-a-mockingbird', 'Pulitzer Prize-winning novel by Harper Lee', 1599, 'USD', (SELECT id FROM categories WHERE slug = 'books'), ARRAY[]::TEXT[], 120, true),
  ('1984', 'nineteen-eighty-four', 'Dystopian novel by George Orwell', 1699, 'USD', (SELECT id FROM categories WHERE slug = 'books'), ARRAY[]::TEXT[], 200, true)
ON CONFLICT (slug) DO UPDATE
SET
  name = EXCLUDED.name,
  description = EXCLUDED.description,
  price_cents = EXCLUDED.price_cents,
  currency = EXCLUDED.currency,
  category_id = EXCLUDED.category_id,
  stock_quantity = EXCLUDED.stock_quantity,
  is_active = EXCLUDED.is_active,
  updated_at = NOW();
SQL

echo "✅ Sample data seeded successfully."
echo ""
echo "Verify with:"
echo "  curl http://localhost:8080/api/v1/categories"
echo "  curl http://localhost:8080/api/v1/products"
