# Build stage
FROM golang:1.26-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod files
COPY services/order/go.mod services/order/go.sum* ./services/order/
COPY proto/ ./proto/

# Download dependencies
WORKDIR /app/services/order
RUN go mod download

# Copy source code
WORKDIR /app
COPY services/order/ ./services/order/

# Build the application
WORKDIR /app/services/order
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /order-service ./cmd/order

# Runtime stage
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Create non-root user
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

WORKDIR /app

# Copy binary from builder
COPY --from=builder /order-service .

# Change ownership
RUN chown -R appuser:appuser /app

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 50055

# Run the binary
CMD ["./order-service"]
