# Architecture Documentation

## Overview

This project implements a RESTful Product Catalog API written in Go.

The application uses:

- Gorilla Mux for HTTP routing
- A layered microservice architecture
- Two different storage implementations:
  - `MemoryStore`
  - `PostgresStore`

The application starts in `cmd/api/main.go`.

At startup, the application checks whether the environment variable `DB_HOST` is configured.

- If `DB_HOST` exists, the application connects to PostgreSQL using `PostgresStore`
- Otherwise, it falls back to the in-memory implementation `MemoryStore`

This allows the application to run both:
- locally without a database
- and in a persistent Docker/PostgreSQL environment

---

# Project Structure

## cmd/api

Contains the application entry point:

```text
cmd/api/main.go
```

Responsibilities:
- create router
- configure storage
- start HTTP server
- load environment variables

---

## internal/handler

Contains the HTTP handlers.

Responsibilities:
- process HTTP requests
- parse JSON payloads
- validate request data
- call store layer
- return JSON responses

Important handlers:
- `GetProducts`
- `GetProduct`
- `CreateProduct`
- `UpdateProduct`
- `DeleteProduct`

---

## internal/store

Contains the storage implementations.

### MemoryStore

Stores products in memory using a Go map.

### PostgresStore

Stores products in a PostgreSQL database using SQL queries.

Responsibilities:
- CRUD operations
- persistence
- database access

---

## internal/model

Contains the product model and validation logic.

Example:

```go
type Product struct {
    ID    int
    Name  string
    Price float64
}
```

Responsibilities:
- product structure
- validation rules

---

# Request Flow

The API follows a layered request flow:

```text
HTTP Request
    ↓
Gorilla Mux Router
    ↓
Handler
    ↓
Store Layer
    ↓
MemoryStore or PostgreSQL Database
    ↓
JSON Response
```

---

# Detailed Request Flow

## Example: GET /products

### Step 1 — Client Request

A client sends:

```http
GET /products
```

---

### Step 2 — Router

Gorilla Mux receives the request.

The route:

```go
r.HandleFunc("/products", h.GetProducts).Methods("GET")
```

matches the request.

---

### Step 3 — Handler

The `GetProducts` handler is executed:

```go
func (h *Handler) GetProducts(w http.ResponseWriter, r *http.Request)
```

The handler:
- processes the request
- calls the store layer
- prepares the response

---

### Step 4 — Store Layer

The handler calls:

```go
h.Store.GetAll()
```

Depending on the configuration:

- `MemoryStore.GetAll()`
- or `PostgresStore.GetAll()`

is executed.

---

### Step 5 — Data Access

### MemoryStore

Products are loaded from:

```go
map[int]model.Product
```

### PostgresStore

Products are loaded using SQL:

```sql
SELECT id, name, price FROM products
```

---

### Step 6 — JSON Response

The handler converts the result into JSON:

```json
[
  {
    "id": 1,
    "name": "Keyboard",
    "price": 49.99
  }
]
```

The API returns:

```http
HTTP 200 OK
```

---

# Architecture Diagram

```text
+-------------------+
| HTTP Client       |
+-------------------+
          |
          v
+-------------------+
| Gorilla Mux       |
| Router            |
+-------------------+
          |
          v
+-------------------+
| HTTP Handler      |
+-------------------+
          |
          v
+-------------------+
| Store Layer       |
+-------------------+
      /       \
     /         \
    v           v
+-----------+  +----------------+
| Memory    |  | PostgreSQL     |
| Store     |  | Database       |
+-----------+  +----------------+
```

---

# API Endpoints

| Method | Endpoint | Description |
|---|---|---|
| GET | `/health` | Health check |
| GET | `/products` | List all products |
| POST | `/products` | Create product |
| GET | `/products/{id}` | Get product by ID |
| PUT | `/products/{id}` | Update product |
| DELETE | `/products/{id}` | Delete product |

---

# MemoryStore

## Description

`MemoryStore` stores all products in application memory.

It uses:

```go
map[int]model.Product
```

and protects concurrent access using:

```go
sync.RWMutex
```

---

## When to Use MemoryStore

MemoryStore is useful for:

- local development
- unit testing
- quick prototypes
- simple demos
- environments without a database

---

## Advantages

- very simple setup
- no external dependencies
- fast read/write operations
- ideal for automated tests
- easy debugging

---

## Disadvantages

- data is lost after restart
- no persistence
- not suitable for production
- limited scalability
- cannot share state between instances

---

# PostgresStore

## Description

`PostgresStore` stores products in a PostgreSQL database.

It uses SQL queries for CRUD operations.

Example:

```sql
SELECT id, name, price FROM products
```

The database table is created automatically:

```sql
CREATE TABLE IF NOT EXISTS products (
    id    SERIAL PRIMARY KEY,
    name  TEXT NOT NULL,
    price NUMERIC(10,2) NOT NULL DEFAULT 0
)
```

---

## When to Use PostgresStore

PostgresStore is useful for:

- Docker Compose environments
- integration testing
- persistent storage
- production-like systems
- scalable deployments

---

## Advantages

- persistent storage
- survives application restarts
- suitable for production
- supports multiple instances
- realistic database environment
- better scalability

---

## Disadvantages

- requires PostgreSQL
- more configuration
- slower than in-memory access
- more complex setup
- database/network dependency

---

# Comparison: MemoryStore vs PostgresStore

| Aspect | MemoryStore | PostgresStore |
|---|---|---|
| Storage | Go map in memory | PostgreSQL database |
| Persistence | No | Yes |
| Setup Complexity | Very low | Higher |
| Performance | Very fast | Database dependent |
| Production Usage | No | Yes |
| Scalability | Limited | Better |
| External Dependency | None | PostgreSQL |
| Best Use Case | Unit tests, local dev | Persistent deployments |

---

# Summary

The application follows a layered microservice architecture.

The request flow is:

```text
HTTP Request → Router → Handler → Store → Database
```

The router handles incoming requests and forwards them to the correct handler.

Handlers process requests and responses, while the store layer handles data access.

The architecture cleanly separates:
- HTTP logic
- business logic
- persistence logic

`MemoryStore` is ideal for development and testing because it is simple and fast.

`PostgresStore` is better suited for persistent and production-like environments because data survives application restarts and can scale better.

