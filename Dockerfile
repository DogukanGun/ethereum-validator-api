# Build stage
FROM golang:1.24-alpine AS builder

# Set working directory
WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application with proper flags
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o ethereum-validator-api .

# Final stage
FROM alpine:3.19

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy binary from build stage
COPY --from=builder /app/ethereum-validator-api .
COPY --from=builder /app/.env .

# Expose the port the app runs on
EXPOSE 3001

# Run with proper user permissions
RUN adduser -D -u 1000 appuser
USER appuser

# Run the binary
CMD ["./ethereum-validator-api"] 