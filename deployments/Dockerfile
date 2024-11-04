# Build stage
FROM golang:1.21-alpine AS builder

# Install necessary build tools
RUN apk add --no-cache git make protoc

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN make build

# Final stage
FROM alpine:3.18

WORKDIR /app

# Install necessary runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Copy the binary from builder
COPY --from=builder /app/bin/server /app/server
COPY --from=builder /app/configs /app/configs

# Create a non-root user
RUN adduser -D -g '' appuser
USER appuser

# Expose gRPC and HTTP ports
EXPOSE 8080 9090

# Run the service
CMD ["/app/server"]