# Build stage
FROM rust:1-alpine AS builder

# Install build dependencies
RUN apk add --no-cache musl-dev protobuf-dev

# Set working directory
WORKDIR /app

# Copy proto files
COPY proto/ /app/proto/

# Copy Cargo files
COPY services/payment/Cargo.toml services/payment/Cargo.lock services/payment/build.rs /app/services/payment/

# Build from the service directory so build.rs relative proto paths resolve.
WORKDIR /app/services/payment

# Create dummy main to cache dependencies
RUN mkdir src && echo "fn main() {}" > src/main.rs
RUN cargo build --release --locked
RUN rm -rf src

# Copy source code
COPY services/payment/src /app/services/payment/src

# Build the application
RUN touch src/main.rs && cargo build --release --locked

# Runtime stage
FROM alpine:3.23

# Install runtime dependencies
RUN apk --no-cache add ca-certificates libgcc

# Create non-root user
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/services/payment/target/release/payment-service .

# Change ownership
RUN chown -R appuser:appuser /app

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 50056

# Run the binary
CMD ["./payment-service"]
