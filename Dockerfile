FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN go build -o server ./cmd/server

# Create final lightweight image
FROM alpine:latest

WORKDIR /app

# Install CA certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Copy binary from builder
COPY --from=builder /app/server /app/server
COPY --from=builder /app/internal/templates /app/internal/templates

# Create volume mount point for PocketBase data
VOLUME /app/pb_data

# Environment variables with defaults
ENV PORT=8080
ENV PB_DATA_DIR=/app/pb_data

# Expose the port
EXPOSE 8080

# Run the server
CMD ["/app/server"]