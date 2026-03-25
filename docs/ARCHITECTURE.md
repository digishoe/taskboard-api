# Architecture

## Overview

The taskboard API follows a three-layer architecture:

```
HTTP Request
    │
    ▼
┌──────────────────┐
│  CORS Middleware  │  Sets Access-Control-* headers, handles OPTIONS preflight
└────────┬─────────┘
         │
         ▼
┌──────────────────┐
│    net/http Mux   │  Go 1.22+ pattern routing (method + path)
└────────┬─────────┘
         │
         ▼
┌──────────────────┐
│     Handlers      │  Parse request, call store, write JSON response
│  handler/*.go     │
└────────┬─────────┘
         │
         ▼
┌──────────────────┐
│      Store        │  SQL queries, transactions
│  store/sqlite.go  │
└────────┬─────────┘
         │
         ▼
┌──────────────────┐
│     SQLite DB     │  Via modernc.org/sqlite (pure Go, no CGO)
└──────────────────┘
```

## Code Map

- **`main.go`** — Application entrypoint. Reads `PORT` and `DB_PATH` env vars, creates the store, wires up handlers to routes, starts the HTTP server with CORS middleware.

- **`internal/store/sqlite.go`** — Database layer. `New()` opens the database, runs migrations (CREATE TABLE IF NOT EXISTS), and seeds a default board with three columns on first run. Exposes methods for all CRUD operations.

- **`internal/model/models.go`** — Data types: `Board`, `Column`, `Task`. Used for both database scanning and JSON serialization. Nested structs use `omitempty` to keep list responses clean.

- **`internal/handler/boards.go`** — Board endpoints: list, create, get (with nested columns and tasks).

- **`internal/handler/columns.go`** — Column endpoint: create column within a board.

- **`internal/handler/tasks.go`** — Task endpoints: create, update (move between columns), delete.

## Design Decisions

### Pure Go SQLite
Using `modernc.org/sqlite` instead of `mattn/go-sqlite3` eliminates the CGO dependency. This simplifies cross-compilation and CI pipelines — no C compiler needed.

### Stdlib Router
Go 1.22 introduced method-aware pattern routing in `net/http.ServeMux` (e.g., `GET /api/boards/{id}`). This eliminates the need for third-party routers like chi or gorilla/mux for simple APIs.

### In-Memory Test Databases
Every test creates a fresh `:memory:` SQLite database. This provides complete isolation between tests with no cleanup needed, and tests run fast since there is no disk I/O.

### Seed Data
On first run (empty database), a "Default Board" is created with three standard kanban columns: To Do, In Progress, Done. This means the API is immediately usable without any setup.

### Flat Store Interface
All database methods live on `*store.SQLiteStore` rather than behind an interface. For a small project this keeps things simple. If mocking is needed later, extract an interface at that point.

### CORS Middleware
A simple middleware adds permissive CORS headers (`Access-Control-Allow-Origin: *`) for local frontend development. In production, this should be tightened to specific origins.
