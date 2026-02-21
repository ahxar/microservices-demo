# Build stage
FROM golang:1.26-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy everything
COPY proto/ ./proto/
COPY services/user/ ./services/user/

# Build the application (with -mod=mod to allow dependency updates)
WORKDIR /app/services/user
RUN CGO_ENABLED=0 GOOS=linux go build -mod=mod -o /user-service ./cmd/user

# Runtime stage
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Create non-root user
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

WORKDIR /app

# Copy binary from builder
COPY --from=builder /user-service .

# Change ownership
RUN chown -R appuser:appuser /app

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 50051

# Run the binary
CMD ["./user-service"]
