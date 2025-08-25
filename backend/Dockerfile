# Use the official Go image as the base image
FROM golang:1.23.2-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Install git (needed for Go modules)
RUN apk add --no-cache git

# Copy go mod and sum files
COPY app/go.mod app/go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY app/ .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/app

# Start a new stage from scratch
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Create app directory
WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Copy migrations
COPY --from=builder /app/db ./db

# Copy .env file (optional, can be overridden by docker-compose)
COPY --from=builder /app/.env .

# Expose port
EXPOSE 8080

# Command to run
CMD ["./main"]