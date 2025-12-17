# Docker Support

[English](docker.md) | [日本語](docker_ja.md)

You can use goreload within Docker containers to enable hot reloading in your containerized development environment.

## Dockerfile Setup

To use goreload in your project, install it using `go install` in your `Dockerfile`.

```dockerfile
FROM golang:1.25

# Install goreload
RUN go install github.com/taro33333/goreload/cmd/goreload@latest

WORKDIR /app

# Copy go module files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Set goreload as the entry point
CMD ["goreload"]
```

## Docker Compose Setup

For hot reloading to work, you **must** mount your local source code directory into the container. This allows goreload to detect file changes made on your host machine.

```yaml
services:
  app:
    build: .
    # Mount the current directory to /app in the container
    volumes:
      - .:/app
    # Expose ports as needed
    ports:
      - "8080:8080"
```

## Configuration Tips

### Build Output

By default, goreload builds the binary to `./tmp/main`. Ensure your container user has write permissions to this directory.

### Linux/Docker File Watching

goreload uses `fsnotify`, which relies on `inotify` on Linux (and thus inside Docker containers). This works seamlessly with Docker for Mac/Windows file sharing.

### Multi-stage Builds

For production builds, you typically don't need goreload. You can use multi-stage builds to separate the development environment from the production image.

```dockerfile
# Development stage
FROM golang:1.25 AS dev
RUN go install github.com/taro33333/goreload/cmd/goreload@latest
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
CMD ["goreload"]

# Production stage
FROM golang:1.25 AS builder
WORKDIR /app
COPY . .
RUN go build -o main .

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/main .
CMD ["./main"]
```
