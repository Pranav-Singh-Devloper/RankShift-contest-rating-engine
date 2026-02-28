# backend/Dockerfile
# Build Stage
FROM golang:1.23-alpine AS builder

# Install required system dependencies for Prisma
RUN apk add --no-cache openssl git ca-certificates

WORKDIR /app

# Cache Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application
COPY . .

# Generate the Prisma Go Client
RUN go run github.com/steebchen/prisma-client-go generate

# Build the Go binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/engine-api ./cmd/api

# Final Production Stage
FROM alpine:latest

# Install OpenSSL (required by Prisma engine at runtime)
RUN apk add --no-cache openssl ca-certificates

WORKDIR /app

# Copy the compiled binary and the Prisma schema from the builder stage
COPY --from=builder /app/engine-api .
COPY --from=builder /app/prisma ./prisma

# Expose the port
EXPOSE 8080

# Run the binary
CMD ["./engine-api"]