# Build stage
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Copy proto files
COPY proto /app/proto

# Copy notification service files
COPY services/notification /app/services/notification

# Build the service
WORKDIR /app/services/notification
RUN go mod download && \
    CGO_ENABLED=0 GOOS=linux go build -o /app/notification-service ./cmd/notification

# Runtime stage
FROM alpine:latest

WORKDIR /app

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Copy binary from builder
COPY --from=builder /app/notification-service /app/notification-service

# Expose gRPC port
EXPOSE 50057

# Run the service
CMD ["/app/notification-service"]
