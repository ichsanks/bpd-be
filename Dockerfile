# Specify Go version as build argument
ARG GO_VERSION=1.18

# First stage: Builder
FROM golang:${GO_VERSION}-alpine AS builder

# Install build dependencies
RUN apk update && \
    apk add --no-cache \
    git \
    make \
    build-base

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum first to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Ensure .env file exists (copy from example if it doesn't)
COPY .env .env

# Generate any required code
RUN go generate ./...

# Build the binary with optimizations
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -o goBinary .

# Second stage: Final image
FROM alpine:latest

# Install necessary runtime dependencies
RUN apk update && \
    apk add --no-cache \
    ca-certificates \
    tzdata

# Set timezone
ENV TZ=Asia/Jakarta

# Create non-root user
RUN adduser -D appuser

# Set working directory
WORKDIR /app

# Create upload directory structure with proper permissions
RUN mkdir -p /app/uploads && \
    mkdir -p /app/temp && \
    chown -R appuser:appuser /app/uploads && \
    chown -R appuser:appuser /app/temp && \
    chmod 755 /app/uploads && \
    chmod 755 /app/temp 

# Copy binary and env files
COPY --from=builder /app/goBinary .
COPY --from=builder /app/.env .

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8076

# Run the entrypoint script
CMD ["./goBinary"]