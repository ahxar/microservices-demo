#!/bin/bash

# Script to apply all database migrations
# Run this after starting docker compose

set -e

POSTGRES_CONTAINER="microservices-postgres"
POSTGRES_USER="postgres"

echo "Applying database migrations..."
echo ""

# Function to apply migrations for a service
apply_migrations() {
    local service=$1
    local db_name=$2

    echo "📦 Applying migrations for $service service..."

    if [ -d "services/$service/migrations" ]; then
        for migration in services/$service/migrations/*.up.sql; do
            if [ -f "$migration" ]; then
                echo "  ✓ Applying $(basename $migration)"
                docker exec -i $POSTGRES_CONTAINER psql -U $POSTGRES_USER -d $db_name < "$migration"
            fi
        done
        echo "  ✅ $service migrations complete"
    else
        echo "  ⚠️  No migrations found for $service"
    fi
    echo ""
}

# Check if postgres container is running
if ! docker ps | grep -q $POSTGRES_CONTAINER; then
    echo "❌ Error: PostgreSQL container is not running"
    echo "Please start docker compose first:"
    echo "  docker compose -f deployments/docker/docker-compose.yml up -d postgres"
    exit 1
fi

# Wait for postgres to be ready
echo "⏳ Waiting for PostgreSQL to be ready..."
sleep 3

# Apply migrations for each service
apply_migrations "user" "user_db"
apply_migrations "catalog" "catalog_db"
apply_migrations "cart" "cart_db"
apply_migrations "order" "order_db"
apply_migrations "payment" "payment_db"
apply_migrations "shipping" "shipping_db"
apply_migrations "notification" "notification_db"

echo "🎉 All migrations applied successfully!"
echo ""
echo "You can now start all services:"
echo "  docker compose -f deployments/docker/docker-compose.yml up -d"
