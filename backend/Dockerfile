# backend/Dockerfile
FROM golang:1.23-alpine AS builder

# 1. Install system dependencies
RUN apk add --no-cache openssl git ca-certificates

WORKDIR /app

# 2. Debugging Step: List files to ensure we are in the right place
# This will show up in your Render build logs
RUN ls -la 

# 3. Copy dependency files
# If Render Root Directory is 'backend', these are in the current folder
COPY go.mod go.sum ./

# 4. Explicitly set proxy to bypass DNS issues in some build environments
ENV GOPROXY=https://proxy.golang.org,direct
RUN go mod download

# 5. Copy the rest of the backend source code
COPY . .

# 6. Generate Prisma and Build
# Ensure your package name in main.go matches the module name in go.mod
RUN go run github.com/steebchen/prisma-client-go generate
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/engine-api ./cmd/api

# Final Production Stage
FROM alpine:latest
RUN apk add --no-cache openssl ca-certificates
WORKDIR /app

# Copy the binary and prisma folder for the runtime engine
COPY --from=builder /app/engine-api .
COPY --from=builder /app/prisma ./prisma

EXPOSE 8080
CMD ["./engine-api"]