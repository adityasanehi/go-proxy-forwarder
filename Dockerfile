FROM golang:1.22-alpine

# Install SQLite
RUN apk add --no-cache sqlite

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN go build -o /go-proxy-rotator

# Create directory for database
RUN mkdir -p /data

EXPOSE 3000

# Set environment variables
ENV DATABASE_PATH=/data/proxies.db

CMD [ "/go-proxy-rotator" ]
