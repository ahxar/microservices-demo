# Build stage
FROM golang:1.26-alpine AS builder
RUN apk add --no-cache git
WORKDIR /app
COPY proto/ ./proto/
COPY gateway/ ./gateway/
WORKDIR /app/gateway
RUN CGO_ENABLED=0 GOOS=linux go build -mod=mod -o /gateway ./cmd/gateway

# Runtime stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
RUN addgroup -g 1000 appuser && adduser -D -u 1000 -G appuser appuser
WORKDIR /app
COPY --from=builder /gateway .
RUN chown -R appuser:appuser /app
USER appuser
EXPOSE 8080
CMD ["./gateway"]
