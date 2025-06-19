# Build stage
FROM golang:1.23-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the applications
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/api
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o worker ./cmd/worker

# API stage (default)
FROM alpine:latest AS api

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the API binary from builder
COPY --from=builder /app/main .

# Copy .env.example as .env (can be overridden with volume mount)
COPY --from=builder /app/.env.example .env

# Expose port
EXPOSE 8080

# Run the API binary
CMD ["./main"]

# Worker stage
FROM alpine:latest AS worker

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the worker binary from builder
COPY --from=builder /app/worker .

# Copy .env.example as .env (can be overridden with volume mount)
COPY --from=builder /app/.env.example .env

# Run the worker binary
CMD ["./worker"]