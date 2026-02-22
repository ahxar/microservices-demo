#!/bin/bash

# Seed sample catalog data directly into PostgreSQL.
# This script is idempotent and can be re-run safely.

set -euo pipefail

POSTGRES_CONTAINER="microservices-postgres"
POSTGRES_USER="postgres"
CATALOG_DB="catalog_db"
USER_DB="user_db"
ADMIN_EMAIL="admin@example.com"
ADMIN_PASSWORD_BCRYPT='$2a$10$BGrO52fnYe5vwaKYQ1QACOZI89hpbOrE5FawM2AWa0DFbL8q6ofMa'

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
  ('Wireless Headphones', 'wireless-headphones', 'Premium noise-canceling wireless headphones with 30-hour battery life', 29900, 'USD', (SELECT id FROM categories WHERE slug = 'electronics'), ARRAY['http://localhost:3000/images/products/wireless-headphones.jpg']::TEXT[], 50, true),
  ('Smart Watch', 'smart-watch', 'Fitness tracking smartwatch with heart rate monitor and GPS', 39900, 'USD', (SELECT id FROM categories WHERE slug = 'electronics'), ARRAY['http://localhost:3000/images/products/smart-watch.jpg']::TEXT[], 30, true),
  ('Laptop Backpack', 'laptop-backpack', 'Durable backpack with padded laptop compartment, fits up to 15 inch laptops', 5900, 'USD', (SELECT id FROM categories WHERE slug = 'electronics'), ARRAY['http://localhost:3000/images/products/laptop-backpack.jpg']::TEXT[], 100, true),
  ('Cotton T-Shirt', 'cotton-t-shirt', 'Comfortable 100% cotton t-shirt, available in multiple colors', 1999, 'USD', (SELECT id FROM categories WHERE slug = 'clothing'), ARRAY['http://localhost:3000/images/products/cotton-t-shirt.png']::TEXT[], 200, true),
  ('Denim Jeans', 'denim-jeans', 'Classic fit denim jeans, durable and stylish', 4999, 'USD', (SELECT id FROM categories WHERE slug = 'clothing'), ARRAY['http://localhost:3000/images/products/denim-jeans.jpg']::TEXT[], 75, true),
  ('Winter Jacket', 'winter-jacket', 'Warm and waterproof winter jacket with hood', 12900, 'USD', (SELECT id FROM categories WHERE slug = 'clothing'), ARRAY['http://localhost:3000/images/products/winter-jacket.jpg']::TEXT[], 40, true),
  ('The Great Gatsby', 'the-great-gatsby', 'Classic American novel by F. Scott Fitzgerald', 1499, 'USD', (SELECT id FROM categories WHERE slug = 'books'), ARRAY['http://localhost:3000/images/products/the-great-gatsby.jpg']::TEXT[], 150, true),
  ('To Kill a Mockingbird', 'to-kill-a-mockingbird', 'Pulitzer Prize-winning novel by Harper Lee', 1599, 'USD', (SELECT id FROM categories WHERE slug = 'books'), ARRAY['http://localhost:3000/images/products/to-kill-a-mockingbird.jpg']::TEXT[], 120, true),
  ('1984', 'nineteen-eighty-four', 'Dystopian novel by George Orwell', 1699, 'USD', (SELECT id FROM categories WHERE slug = 'books'), ARRAY['http://localhost:3000/images/products/nineteen-eighty-four.jpg']::TEXT[], 200, true)
ON CONFLICT (slug) DO UPDATE
SET
  name = EXCLUDED.name,
  description = EXCLUDED.description,
  price_cents = EXCLUDED.price_cents,
  currency = EXCLUDED.currency,
  category_id = EXCLUDED.category_id,
  image_urls = EXCLUDED.image_urls,
  stock_quantity = EXCLUDED.stock_quantity,
  is_active = EXCLUDED.is_active,
  updated_at = NOW();
SQL

echo "✅ Sample data seeded successfully."

docker exec -i "${POSTGRES_CONTAINER}" psql -U "${POSTGRES_USER}" -d "${USER_DB}" <<SQL
WITH seeded_admin AS (
  INSERT INTO users (email, password_hash, role)
  VALUES ('${ADMIN_EMAIL}', '${ADMIN_PASSWORD_BCRYPT}', 'admin')
  ON CONFLICT (email) DO UPDATE
  SET
    password_hash = EXCLUDED.password_hash,
    role = 'admin',
    updated_at = NOW()
  RETURNING id
)
INSERT INTO profiles (user_id, first_name, last_name, phone, avatar_url)
SELECT id, 'Admin', 'User', '', ''
FROM seeded_admin
ON CONFLICT (user_id) DO UPDATE
SET
  first_name = EXCLUDED.first_name,
  last_name = EXCLUDED.last_name,
  phone = EXCLUDED.phone,
  avatar_url = EXCLUDED.avatar_url;
SQL

echo "✅ Admin user seeded successfully."
echo ""
echo "Verify with:"
echo "  curl http://localhost:8080/api/v1/categories"
echo "  curl http://localhost:8080/api/v1/products"
echo ""
echo "Admin login seed:"
echo "  email: ${ADMIN_EMAIL}"
echo "  password: admin123"
