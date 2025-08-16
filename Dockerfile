# Multi-stage build for API Gateway (main entry point)
FROM golang:1.23.5-alpine AS builder

WORKDIR /app

# Copy API Gateway files
COPY clans/api-gateway/go.mod clans/api-gateway/go.sum ./
RUN go mod download

COPY clans/api-gateway/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/main.go

# Production stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/main .

# Expose API Gateway port
EXPOSE 8000

CMD ["./main"]