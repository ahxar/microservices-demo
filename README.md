# Microservices Demo - Production-Grade E-Commerce Platform

A modern, production-ready microservices architecture e-commerce application inspired by Google Cloud's microservices-demo, built with **Go, Rust, Next.js**, and best practices (KISS, SOLID, YAGNI).

## 🏗️ Architecture

See **[ARCHITECTURE.md](ARCHITECTURE.md)** for detailed Mermaid flowcharts including:

- Complete system architecture diagram
- Request flow sequences (registration, browsing, cart, checkout)
- API Gateway route mapping
- Database schema per service
- Technology stack breakdown

**Quick Overview:**

```
Frontend (Next.js)
    ↓ HTTP/REST
API Gateway (Go) - Port 8080
    ↓ gRPC
┌─────────────────────────────────────┐
│ ✅ User Service (50051)             │
│ ✅ Catalog Service (50052)          │
│ ✅ Cart Service (50053)             │
│ ✅ Order Service (50055)            │
│ ✅ Payment Service (Rust, 50056)    │
│ ✅ Shipping Service (50058)         │
│ ✅ Notification Service (50057)     │
└─────────────────────────────────────┘
    ↓
PostgreSQL (6 DBs) + Redis
```

## 🚀 Tech Stack

| Layer            | Technology                                        |
| ---------------- | ------------------------------------------------- |
| Frontend         | Next.js 16+ (App Router), shadcn/ui, Tailwind CSS |
| API Gateway      | Go (chi router)                                   |
| Backend Services | Go (6 services), Rust (1 service - Payment)       |
| Inter-service    | gRPC with Protocol Buffers                        |
| Databases        | PostgreSQL (per service), Redis (cart/cache)      |
| Deployment       | Docker Compose (local), Kubernetes (prod)         |

## 📦 Services

| Service              | Language | Port  | Database   | Responsibility                      |
| -------------------- | -------- | ----- | ---------- | ----------------------------------- |
| User Service         | Go       | 50051 | PostgreSQL | Auth, profiles, addresses, wishlist |
| Catalog Service      | Go       | 50052 | PostgreSQL | Products, categories, inventory     |
| Cart Service         | Go       | 50053 | Redis      | Shopping cart (ephemeral)           |
| Order Service        | Go       | 50055 | PostgreSQL | Order processing, orchestration     |
| Payment Service      | Rust     | 50056 | PostgreSQL | Payment processing, tokenization    |
| Shipping Service     | Go       | 50058 | PostgreSQL | Shipping quotes, tracking           |
| Notification Service | Go       | 50057 | PostgreSQL | Email notifications                 |
| API Gateway          | Go       | 8080  | Redis      | Auth, routing, rate limiting        |

### Why Payment Service in Rust?

- **Security-critical**: Handles sensitive payment data with memory safety
- **Performance-critical**: Low-latency requirement with no GC pauses
- **Type safety**: Strong type system prevents common vulnerabilities

## 🎯 Current Status

### ✅ Implemented Services (8/8)

- **User Service** (Port 50051): JWT auth, profiles, addresses, wishlists
- **Catalog Service** (Port 50052): Products, categories, inventory, search
- **Cart Service** (Port 50053): Redis-backed shopping cart with TTL
- **Order Service** (Port 50055): Checkout orchestration and order lifecycle
- **Payment Service** (Rust, Port 50056): Payment methods and transactions
- **Shipping Service** (Port 50058): Quotes, shipment creation, tracking
- **Notification Service** (Port 50057): Email templates and SMTP delivery
- **API Gateway** (Port 8080): REST routing, auth, rate limiting, order routes

### 📊 Detailed Progress

See **[TASKS.md](TASKS.md)** for the current completion checklist and remaining polish/testing work.

## 🛠️ Prerequisites

- **Go** 1.21+
- **Rust** 1.75+ (for Payment Service)
- **Node.js** 18+ (for Frontend)
- **Docker** & **Docker Compose**
- **PostgreSQL** 16+ (via Docker)
- **Redis** 7+ (via Docker)
- **Protocol Buffers** compiler (`protoc`)

### Install protoc plugins:

```bash
# Go plugins
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Add to PATH
export PATH="$PATH:$(go env GOPATH)/bin"
```

## 🚀 Quick Start

### 1. Generate Protocol Buffers

```bash
make proto-gen
```

### 2. Start Infrastructure (PostgreSQL, Redis, Mailhog)

```bash
docker compose -f deployments/docker/docker-compose.yml up -d postgres redis mailhog
```

### 3. Run Database Migrations

```bash
make migrate
```

### 4. Seed Sample Data (Optional)

```bash
make seed
```

### 5. Start All Services

```bash
make up
```

### 6. Verify Services

```bash
# Check service health
curl http://localhost:8080/health

# View logs
make logs

# Check Mailhog (email testing)
open http://localhost:8025
```

## 📝 API Documentation

### Authentication

```bash
# Register new user
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123",
    "first_name": "John",
    "last_name": "Doe"
  }'

# Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

### Product Catalog

```bash
# List products
curl http://localhost:8080/api/v1/products?page=1&page_size=10

# Get product by ID
curl http://localhost:8080/api/v1/products/{id}

# Search products
curl "http://localhost:8080/api/v1/products/search?q=laptop"

# List categories
curl http://localhost:8080/api/v1/categories
```

### Shopping Cart (Authenticated)

```bash
# Get cart
curl http://localhost:8080/api/v1/cart \
  -H "Authorization: Bearer {access_token}"

# Add item to cart
curl -X POST http://localhost:8080/api/v1/cart/items \
  -H "Authorization: Bearer {access_token}" \
  -H "Content-Type: application/json" \
  -d '{
    "product_id": "uuid",
    "product_name": "Laptop",
    "quantity": 1,
    "unit_price": {"amount_cents": 99900, "currency": "USD"}
  }'
```

## 🧪 Testing

### gRPC Testing with grpcurl

```bash
# Install grpcurl
brew install grpcurl  # macOS
# or
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest

# List services
grpcurl -plaintext localhost:50051 list

# Test User Service
grpcurl -plaintext \
  -d '{
    "email": "test@example.com",
    "password": "password123",
    "first_name": "Test",
    "last_name": "User"
  }' \
  localhost:50051 user.v1.UserService/Register

# Test Catalog Service
grpcurl -plaintext localhost:50052 catalog.v1.CatalogService/ListProducts
```

### Unit Tests

```bash
make test
```

## 📂 Project Structure

```
/microservices-demo/
├── proto/                      # Shared protobuf definitions
│   ├── common/v1/             # Money, Address, Pagination
│   ├── user/v1/
│   ├── catalog/v1/
│   ├── cart/v1/
│   ├── order/v1/
│   ├── payment/v1/
│   ├── shipping/v1/
│   └── notification/v1/
│
├── gateway/                    # API Gateway (Go)
│   ├── cmd/gateway/
│   ├── internal/
│   │   ├── middleware/        # auth, ratelimit, cors
│   │   ├── handler/           # HTTP handlers
│   │   └── client/            # gRPC clients
│   └── Dockerfile
│
├── services/
│   ├── user/                  # Go service
│   │   ├── cmd/user/
│   │   ├── internal/
│   │   │   ├── config/
│   │   │   ├── server/
│   │   │   ├── service/
│   │   │   └── repository/
│   │   ├── migrations/
│   │   └── Dockerfile
│   ├── catalog/
│   ├── cart/
│   ├── order/
│   ├── shipping/
│   ├── notification/
│   └── payment/               # Rust service
│       ├── src/
│       ├── Cargo.toml
│       └── Dockerfile
│
├── deployments/
│   ├── docker/
│   │   ├── docker-compose.yml
│   │   └── init-db.sql
│   └── kubernetes/
│
├── scripts/
│   ├── proto-gen.sh
│   ├── apply-migrations.sh
│   └── seed-data.sh
│
├── Makefile
└── README.md
```

## 🔧 Development

### Run Individual Services Locally

```bash
# User Service
make dev-user

# Catalog Service
make dev-catalog

# Cart Service
make dev-cart

# API Gateway
make dev-gateway
```

### View Service Logs

```bash
# All services
make logs

# Specific service
docker logs -f user-service
docker logs -f api-gateway
```

## 🐳 Docker Commands

```bash
# Build all images
make build

# Start services
make up

# Stop services
make down

# Clean everything (including volumes)
make clean
```

## 🎯 Order Processing Flow

The Order Service orchestrates a complete checkout flow:

1. **Get Cart** - Retrieve cart from Cart Service
2. **Validate Products** - Check products exist in Catalog
3. **Check Inventory** - Verify sufficient stock
4. **Reserve Inventory** - Lock stock for 15 minutes
5. **Calculate Shipping** - Get quote from Shipping Service
6. **Process Payment** - Charge via Payment Service (with idempotency)
7. **Create Order** - Save order to database
8. **Clear Cart** - Remove items from cart
9. **Send Notification** - Email confirmation to user
10. **Return Order** - Return complete order details

## 📊 Database Schema

Each service has its own database following the database-per-service pattern:

- **user_db**: users, profiles, addresses, wishlists
- **catalog_db**: categories, products, inventory_reservations
- **order_db**: orders, order_items, order_status_history
- **payment_db**: payment_methods, transactions
- **shipping_db**: shipments, tracking_events
- **notification_db**: email_templates, notification_logs

## 🔐 Security Features

- **JWT Authentication**: Secure token-based auth with refresh tokens
- **Password Hashing**: bcrypt with salt
- **Card Tokenization**: Never store raw card numbers
- **Idempotency Keys**: Prevent duplicate payments
- **Rate Limiting**: Redis-based token bucket (100 req/min)
- **CORS**: Configurable origins
- **SQL Injection Protection**: Parameterized queries

## 🚀 Deployment

### Docker Compose (Development)

```bash
make up
```

### Kubernetes (Production)

```bash
# Apply base configuration
kubectl apply -k deployments/kubernetes/base

# Or apply environment-specific overlay
kubectl apply -k deployments/kubernetes/overlays/prod
```

## 📈 Observability

### Mailhog (Email Testing)

Access at: http://localhost:8025

All emails sent by the Notification Service appear in Mailhog.

### Logs

All services use structured JSON logging:

```bash
make logs
```

## 🤝 Contributing

See `IMPLEMENTATION_PLAN.md` for detailed implementation phases and remaining tasks.

## 📄 License

MIT License - See LICENSE file for details

## 🙏 Acknowledgments

Inspired by [Google Cloud's microservices-demo](https://github.com/GoogleCloudPlatform/microservices-demo)

---

**Built with ❤️ using Go, Rust, and modern microservices patterns**
