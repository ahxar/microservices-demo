# Build stage
FROM golang:1.26-alpine AS builder
RUN apk add --no-cache git
WORKDIR /app
COPY proto/ ./proto/
COPY services/cart/ ./services/cart/
WORKDIR /app/services/cart
RUN CGO_ENABLED=0 GOOS=linux go build -mod=mod -o /cart-service ./cmd/cart

# Runtime stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
RUN addgroup -g 1000 appuser && adduser -D -u 1000 -G appuser appuser
WORKDIR /app
COPY --from=builder /cart-service .
RUN chown -R appuser:appuser /app
USER appuser
EXPOSE 50052
CMD ["./cart-service"]
