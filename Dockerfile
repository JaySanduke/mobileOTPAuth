# ===== Build Stage =====
FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

# Copy go.mod and go.sum separately to leverage Docker cache
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the binary from the main package in cmd/server
RUN go build -o main ./cmd/server

# ===== Final Stage =====
FROM alpine:latest

# Set working directory
WORKDIR /app

# Copy only the built binary from the builder stage
COPY --from=builder /app/main .

# Expose port
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
