# Build stage
FROM golang:1.23-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files first for dependency resolution
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the entire source code into the container
COPY . .

# Build the binary for Linux
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o dist/mcp-android-adb-server_linux_amd64 .

# Final stage
FROM alpine:latest

# Set working directory
WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /app/dist/mcp-android-adb-server_linux_amd64 /usr/local/bin/mcp-android-adb-server

# Expose any ports the application requires, if necessary
# EXPOSE <port>

# Set the entrypoint to the compiled binary
ENTRYPOINT ["mcp-android-adb-server"]