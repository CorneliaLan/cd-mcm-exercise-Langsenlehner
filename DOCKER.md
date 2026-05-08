# Docker & Docker Compose Analysis

## Dockerfile Analysis

The application uses a multi-stage Docker build.

---

# Stage 1 — Build Stage

```dockerfile
FROM golang:1.26-alpine AS builder
```

This stage uses a Go Alpine image to compile the application.

```dockerfile
WORKDIR /app
```

Sets the working directory inside the container.

```dockerfile
COPY go.mod go.sum ./
RUN go mod download
```

Copies dependency files and downloads all Go modules.

```dockerfile
COPY . .
```

Copies the complete source code into the container.

```dockerfile
RUN CGO_ENABLED=0 GOOS=linux go build -o /api-server ./cmd/api
```

Builds the Go application as a static Linux binary.

### What does `CGO_ENABLED=0` do?

`CGO_ENABLED=0` disables CGO support and creates a statically linked binary.

Advantages:
- smaller image
- fewer dependencies
- easier deployment
- works well in minimal Alpine containers

Without this option, the application could depend on external system libraries.

---

# Stage 2 — Runtime Stage

```dockerfile
FROM alpine:3.19
```

Uses a very small Alpine Linux runtime image.

```dockerfile
RUN apk --no-cache add ca-certificates
```

Installs SSL certificates for secure HTTPS communication.

```dockerfile
WORKDIR /app
COPY --from=builder /api-server .
```

Copies the compiled binary from the builder stage.

```dockerfile
EXPOSE 8080
```

Exposes the API port.

```dockerfile
ENTRYPOINT ["./api-server"]
```

Starts the application.

---

# Why Multi-Stage Builds?

Multi-stage builds separate:
- build environment
- runtime environment

Advantages:
- smaller final image
- fewer security risks
- faster deployments
- cleaner containers

The final image only contains:
- Alpine Linux
- compiled binary

No Go compiler or source code is included.

---

# Image Size Comparison

## Multi-Stage Build

Smaller production-ready image.

Advantages:
- less storage
- faster downloads
- faster container startup
- reduced attack surface

## Single-Stage Build

A single-stage build would include:
- Go compiler
- build tools
- source code
- dependencies

This would significantly increase image size.

---

# Docker Compose Analysis

The `docker-compose.yml` file starts:

- PostgreSQL database
- Product Catalog API

---

# PostgreSQL Service

```yaml
db:
  image: postgres:16-alpine
```

Starts a PostgreSQL container.

Environment variables configure:
- username
- password
- database name

```yaml
volumes:
  - pgdata:/var/lib/postgresql/data
```

This volume persists database data between container restarts.

---

# API Service

```yaml
api:
  build: .
```

Builds the API image from the Dockerfile.

Environment variables configure:
- database host
- database port
- credentials

```yaml
depends_on:
  db:
    condition: service_healthy
```

Ensures the database is healthy before the API starts.

---

# CRUD Test Results

## Health Check

```bash
curl http://localhost:8080/health
```

Response:

```json
{"status":"ok"}
```

---

## Create Products

```bash
curl -X POST http://localhost:8080/products \
  -H "Content-Type: application/json" \
  -d '{"name":"Keyboard","price":49.99}'
```

```bash
curl -X POST http://localhost:8080/products \
  -H "Content-Type: application/json" \
  -d '{"name":"Mouse","price":19.99}'
```

```bash
curl -X POST http://localhost:8080/products \
  -H "Content-Type: application/json" \
  -d '{"name":"Monitor","price":199.99}'
```

Responses:

```json
{"id":1,"name":"Keyboard","price":49.99}
{"id":2,"name":"Mouse","price":19.99}
{"id":3,"name":"Monitor","price":199.99}
```

---

## List Products

```bash
curl http://localhost:8080/products
```

Response:

```json
[
  {"id":1,"name":"Keyboard","price":49.99},
  {"id":2,"name":"Mouse","price":19.99},
  {"id":3,"name":"Monitor","price":199.99}
]
```

---

## Update Product

```bash
curl -X PUT http://localhost:8080/products/2 \
  -H "Content-Type: application/json" \
  -d '{"name":"Gaming Mouse","price":29.99}'
```

The product was updated successfully.

---

## Delete Product

```bash
curl -X DELETE http://localhost:8080/products/1
```

The product was deleted successfully.

---

# Persistence Test

Containers were stopped and restarted using:

```bash
docker compose down
docker compose up
```

After restarting the containers, the remaining products still existed in the database.

This confirms that PostgreSQL persistence works correctly through the Docker volume.

