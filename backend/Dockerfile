# backend/Dockerfile
FROM golang:1.23-alpine AS builder

# 1. Install system dependencies
RUN apk add --no-cache openssl git ca-certificates

WORKDIR /app

# 2. Copy dependency files FIRST (Optimizes build caching)
# Since your Root Directory is /backend, these files are in the current context
COPY go.mod go.sum ./

# 3. Explicitly set proxy and download
ENV GOPROXY=https://proxy.golang.org,direct
RUN go mod download

# 4. Copy the rest of the backend source code
COPY . .

# 5. Generate Prisma and Build
RUN go run github.com/steebchen/prisma-client-go generate
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/engine-api ./cmd/api

# Final Stage
FROM alpine:latest
RUN apk add --no-cache openssl ca-certificates
WORKDIR /app

COPY --from=builder /app/engine-api .
COPY --from=builder /app/prisma ./prisma

EXPOSE 8080
CMD ["./engine-api"]