# Taskboard API — Agent Orientation

## What This Is
A kanban-style taskboard REST API written in Go. Boards contain columns, columns contain tasks. SQLite storage via `modernc.org/sqlite` (pure Go, no CGO).

## Project Structure
```
main.go                    — Entrypoint: router setup, CORS middleware, env vars
internal/
  store/sqlite.go          — Database init, migrations, seed data
  store/sqlite_test.go     — Store-layer tests (in-memory SQLite)
  model/models.go          — Board, Column, Task types
  model/models_test.go     — JSON serialization tests
  handler/boards.go        — GET/POST /api/boards, GET /api/boards/{id}
  handler/boards_test.go   — Board handler tests
  handler/columns.go       — POST /api/boards/{id}/columns
  handler/columns_test.go  — Column handler tests
  handler/tasks.go         — POST /api/columns/{id}/tasks, PUT/DELETE /api/tasks/{id}
  handler/tasks_test.go    — Task handler tests
```

## Running Tests
```bash
go test ./...
```
All tests use in-memory SQLite (`:memory:`). Each test function creates a fresh database — no shared state between tests.

## Running the Server
```bash
go run .                    # default: port 8080, file taskboard.db
PORT=3000 DB_PATH=dev.db go run .
```

## API Endpoints
| Method | Path                        | Description                     |
|--------|-----------------------------|---------------------------------|
| GET    | /healthz                    | Health check (200 OK)           |
| GET    | /api/boards                 | List all boards                 |
| POST   | /api/boards                 | Create board                    |
| GET    | /api/boards/{id}            | Get board with columns & tasks  |
| POST   | /api/boards/{id}/columns    | Create column in board          |
| POST   | /api/columns/{id}/tasks     | Create task in column           |
| PUT    | /api/tasks/{id}             | Update task                     |
| DELETE | /api/tasks/{id}             | Delete task                     |

## Key Design Decisions
- Uses Go 1.22+ `net/http` stdlib router with method+pattern syntax (e.g., `GET /api/boards/{id}`)
- Store layer is the single point of DB access; handlers depend on `*store.SQLiteStore`
- CORS middleware wraps the entire mux for frontend dev
- Seed data: one "Default Board" with columns "To Do", "In Progress", "Done"

## Data Flow
Request → CORS middleware → ServeMux → Handler → Store → SQLite → Response (JSON)

## Configuration
- `PORT` — listen port (default `8080`)
- `DB_PATH` — SQLite file path (default `taskboard.db`)
