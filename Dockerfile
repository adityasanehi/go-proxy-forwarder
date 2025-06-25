# Build stage
FROM golang:1.22-alpine AS builder

# Install build dependencies
RUN apk add --no-cache \
    gcc \
    musl-dev \
    sqlite-dev

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application with CGO enabled
ENV CGO_ENABLED=1
RUN go build -o /go-proxy-rotator

# Runtime stage
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache \
    sqlite \
    ca-certificates

# Create directory for database
RUN mkdir -p /data

# Copy the binary from builder stage
COPY --from=builder /go-proxy-rotator /go-proxy-rotator

# Copy static files and docs
COPY static ./static
COPY docs ./docs

EXPOSE 3000

# Set environment variables
ENV DATABASE_PATH=/data/proxies.db

CMD [ "/go-proxy-rotator" ]
