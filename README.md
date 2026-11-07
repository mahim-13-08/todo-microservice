# Todo Microservice (Go)

A lightweight, production-ready RESTful Todo microservice built with Go's standard library + SQLite.

## Table of Contents

- [Features](#features)
- [Project Structure](#project-structure)
- [Prerequisites](#prerequisites)
- [Getting Started](#getting-started)
- [Configuration](#configuration)
- [API Reference](#api-reference)
- [Error Handling](#error-handling)
- [Running Tests](#running-tests)

---

## Features

- Full CRUD for Todos (Create / Read / Update / Delete)
- SQLite persistence with automatic schema migration
- Clean layered architecture: Handler → Service → Repository
- Request validation with meaningful error messages
- Structured JSON responses with consistent envelope
- Request logging and CORS middleware
- Graceful shutdown on SIGINT / SIGTERM
- Zero external web-framework dependency (uses `net/http` only)

---

## Project Structure

```
todo-microservice/
├── cmd/
│   └── server/
│       └── main.go           # Entry point, wires dependencies
├── internal/
│   ├── handler/
│   │   └── todo_handler.go   # HTTP handlers (transport layer)
│   ├── service/
│   │   └── todo_service.go   # Business logic
│   ├── repository/
│   │   └── todo_repository.go# Data access (SQLite)
│   ├── model/
│   │   └── todo.go           # Domain model & request DTOs
│   └── middleware/
│       └── middleware.go     # Logger, CORS helpers
├── pkg/
│   ├── response/
│   │   └── response.go       # JSON response helpers
│   └── validator/
│       └── validator.go      # Input validation helpers
├── go.mod
├── go.sum
├── Makefile
├── .gitignore
└── README.md
```

---

## Prerequisites

| Tool | Version |
|------|---------|
| Go   | ≥ 1.22  |
| GCC  | Any (required by `go-sqlite3` CGO) |

> **Windows users:** Install [TDM-GCC](https://jmeubank.github.io/tdm-gcc/) or use WSL.

---

## Getting Started

```bash
# 1. Clone the repository
git clone https://github.com/yourusername/todo-microservice.git
cd todo-microservice

# 2. Download dependencies
go mod download

# 3. Build & run
go run ./cmd/server
# or
make run
```

The server starts on **http://localhost:8080** by default.

---

## Configuration

| Environment Variable | Default       | Description             |
|----------------------|---------------|-------------------------|
| `PORT`               | `8080`        | HTTP listen port        |
| `DATABASE_DSN`       | `./todos.db`  | SQLite file path / DSN  |

```bash
PORT=9090 DATABASE_DSN=/data/todos.db go run ./cmd/server
```

---

## API Reference

### Base URL

```
http://localhost:8080/api/v1
```

### Health Check

```
GET /health
```

**Response `200`**
```json
{"status": "ok"}
```

---

### Todo Object

| Field         | Type     | Description                                 |
|---------------|----------|---------------------------------------------|
| `id`          | integer  | Auto-incremented primary key                |
| `title`       | string   | Short title (1–255 chars, required)         |
| `description` | string   | Optional detail (max 1 000 chars)           |
| `status`      | string   | `pending` \| `in_progress` \| `done`        |
| `created_at`  | ISO 8601 | Creation timestamp (UTC)                    |
| `updated_at`  | ISO 8601 | Last-update timestamp (UTC)                 |

---

### Endpoints

#### Create a Todo

```
POST /api/v1/todos
Content-Type: application/json
```

**Request body**
```json
{
  "title": "Buy groceries",
  "description": "Milk, eggs, bread",
  "status": "pending"
}
```

**Response `201 Created`**
```json
{
  "success": true,
  "message": "todo created successfully",
  "data": {
    "id": 1,
    "title": "Buy groceries",
    "description": "Milk, eggs, bread",
    "status": "pending",
    "created_at": "2024-01-15T10:00:00Z",
    "updated_at": "2024-01-15T10:00:00Z"
  }
}
```

---

#### Retrieve All Todos

```
GET /api/v1/todos
```

**Response `200 OK`**
```json
{
  "success": true,
  "message": "todos retrieved successfully",
  "data": [ { ... }, { ... } ]
}
```

---

#### Retrieve a Single Todo

```
GET /api/v1/todos/{id}
```

**Response `200 OK`**
```json
{
  "success": true,
  "message": "todo retrieved successfully",
  "data": { ... }
}
```

**Response `404 Not Found`**
```json
{
  "success": false,
  "error": "todo not found"
}
```

---

#### Update a Todo (full or partial)

```
PUT   /api/v1/todos/{id}
PATCH /api/v1/todos/{id}
Content-Type: application/json
```

All fields are optional; only fields present in the body are updated.

**Request body**
```json
{
  "status": "done"
}
```

**Response `200 OK`**
```json
{
  "success": true,
  "message": "todo updated successfully",
  "data": { ... }
}
```

---

#### Delete a Todo

```
DELETE /api/v1/todos/{id}
```

**Response `200 OK`**
```json
{
  "success": true,
  "message": "todo deleted successfully"
}
```

---

## Error Handling

All error responses follow the same envelope:

```json
{
  "success": false,
  "error": "<human-readable message>"
}
```

| Status | Meaning                                 |
|--------|-----------------------------------------|
| 400    | Bad request / malformed JSON            |
| 404    | Todo not found                          |
| 405    | Method not allowed                      |
| 422    | Validation error (field-level details)  |
| 500    | Internal server error                   |

---

## Running Tests

```bash
go test ./...
```

---

## Makefile Targets

```bash
make run      # go run ./cmd/server
make build    # go build -o bin/server ./cmd/server
make test     # go test ./...
make lint     # golangci-lint run
make clean    # remove build artifacts
```
