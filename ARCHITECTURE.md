# Architecture Flow - Microservices E-commerce Platform

## System Architecture

```plantuml
@startuml
!define RUNNING #90EE90
!define NOTRUNNING #FFB6C1
!define INFRASTRUCTURE #87CEEB
!define FRONTEND #FFD700

skinparam rectangle {
    BorderColor black
    BorderThickness 2
}

package "Client Layer" {
    actor "👤 User Browser" as USER FRONTEND
}

package "Frontend - Port 3000" {
    rectangle "Next.js 15\nReact App" as NEXT FRONTEND
    rectangle "API Client\nAxios + React Query" as API_CLIENT FRONTEND
}

package "API Gateway - Port 8080" {
    rectangle "🌐 API Gateway\nGo + Chi Router" as GATEWAY RUNNING
    rectangle "🔒 Auth Middleware\nJWT Validation" as AUTH_MW RUNNING
    rectangle "⏱️ Rate Limiter\nRedis Token Bucket" as RATE_MW RUNNING
    rectangle "🔓 CORS Middleware" as CORS_MW RUNNING
}

package "Microservices - gRPC" {
    package "✅ Implemented Services" {
        rectangle "👤 User Service\nPort 50051\nGo" as USER_SVC RUNNING
        rectangle "📦 Catalog Service\nPort 50052\nGo" as CATALOG_SVC RUNNING
        rectangle "🛒 Cart Service\nPort 50053\nGo" as CART_SVC RUNNING
        rectangle "📋 Order Service\nPort 50055\nGo" as ORDER_SVC RUNNING
        rectangle "💳 Payment Service\nPort 50056\nRust" as PAYMENT_SVC RUNNING
        rectangle "🚚 Shipping Service\nPort 50058\nGo" as SHIPPING_SVC RUNNING
        rectangle "📧 Notification Service\nPort 50057\nGo" as NOTIF_SVC RUNNING
    }
}

package "Data Layer" {
    database "🗄️ PostgreSQL\nPort 5432\n6 Databases" as POSTGRES INFRASTRUCTURE
    database "⚡ Redis\nPort 6379\nCache & Cart" as REDIS INFRASTRUCTURE
    rectangle "📬 Mailhog\nPort 8025\nEmail Testing" as MAILHOG INFRASTRUCTURE
}

' Client to Frontend
USER -down-> NEXT : HTTPS
NEXT -down-> API_CLIENT : API Calls

' Frontend to Gateway
API_CLIENT -down-> GATEWAY : HTTP REST

' Gateway Middleware
GATEWAY -down-> CORS_MW
CORS_MW -down-> RATE_MW
RATE_MW -down-> AUTH_MW

' Gateway to Services
AUTH_MW -down-> USER_SVC : gRPC
AUTH_MW -down-> CATALOG_SVC : gRPC
AUTH_MW -down-> CART_SVC : gRPC
AUTH_MW -down-> ORDER_SVC : gRPC

' Service Dependencies
ORDER_SVC ..> CART_SVC : gRPC
ORDER_SVC ..> CATALOG_SVC : gRPC
ORDER_SVC ..> PAYMENT_SVC : gRPC
ORDER_SVC ..> SHIPPING_SVC : gRPC
ORDER_SVC ..> NOTIF_SVC : gRPC

' Database Connections
USER_SVC -down-> POSTGRES : SQL
CATALOG_SVC -down-> POSTGRES : SQL
ORDER_SVC ..> POSTGRES : SQL
PAYMENT_SVC ..> POSTGRES : SQL
SHIPPING_SVC ..> POSTGRES : SQL
NOTIF_SVC ..> POSTGRES : SQL

' Redis Connections
CART_SVC -down-> REDIS : TCP
GATEWAY -down-> REDIS : TCP

' Email
NOTIF_SVC ..> MAILHOG : SMTP

@enduml
```

---

## Service Communication Protocol

```plantuml
@startuml
skinparam componentStyle rectangle

component "Client Browser" as A
component "API Gateway" as B
component "User Service" as C
component "Catalog Service" as D
component "Cart Service" as E
component "Order Service" as F
component "Payment Service" as G
component "Shipping Service" as H
component "Notification Service" as I

A --> B : HTTP/REST\nJSON
B --> C : gRPC\nProtocol Buffers
B --> D : gRPC\nProtocol Buffers
B --> E : gRPC\nProtocol Buffers
B --> F : gRPC\nProtocol Buffers

F ..> C : gRPC\nProtocol Buffers
F ..> D : gRPC\nProtocol Buffers
F ..> E : gRPC\nProtocol Buffers
F ..> G : gRPC\nProtocol Buffers
F ..> H : gRPC\nProtocol Buffers
F ..> I : gRPC\nProtocol Buffers
@enduml
```

---

## Legend

- ✅ **Solid Lines**: Active/implemented connections
- ❌ **Dashed Lines**: Future/optional extensions
- 🟢 **Green**: Currently running services
- 🔴 **Pink**: Services not running
- 🔵 **Blue**: Infrastructure components
- 🟡 **Yellow**: Frontend components

---

**Last Updated:** 2026-02-16
**Services Implemented:** 8/8 + frontend
**Status:** Development - integration complete, polish/testing pending
