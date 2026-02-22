# Build stage
FROM golang:1-alpine AS builder
RUN apk add --no-cache git
WORKDIR /app
COPY proto/ ./proto/
COPY services/catalog/ ./services/catalog/
WORKDIR /app/services/catalog
RUN CGO_ENABLED=0 GOOS=linux go build -mod=mod -o /catalog-service ./cmd/catalog

# Runtime stage
FROM alpine:3.23
RUN apk --no-cache add ca-certificates
RUN addgroup -g 1000 appuser && adduser -D -u 1000 -G appuser appuser
WORKDIR /app
COPY --from=builder /catalog-service .
RUN chown -R appuser:appuser /app
USER appuser
EXPOSE 50052
CMD ["./catalog-service"]
