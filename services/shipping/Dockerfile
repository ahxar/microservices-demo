# Build stage
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Copy proto files
COPY proto /app/proto

# Copy shipping service files
COPY services/shipping /app/services/shipping

# Build the service
WORKDIR /app/services/shipping
RUN go mod download && \
    CGO_ENABLED=0 GOOS=linux go build -o /app/shipping-service ./cmd/shipping

# Runtime stage
FROM alpine:latest

WORKDIR /app

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Copy binary from builder
COPY --from=builder /app/shipping-service /app/shipping-service

# Expose gRPC port
EXPOSE 50058

# Run the service
CMD ["/app/shipping-service"]
