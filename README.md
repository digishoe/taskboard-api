# Taskboard API

A kanban-style taskboard REST API built with Go, using SQLite for storage. Boards contain columns, columns contain tasks — drag-and-drop ready.

## Quick Start

```bash
make dev
# Server starts on http://localhost:8080

curl http://localhost:8080/api/boards
```

A default board with three columns (To Do, In Progress, Done) is created automatically on first run.

## API Endpoints

| Method   | Path                        | Body                                                        | Description                    |
|----------|-----------------------------|-------------------------------------------------------------|--------------------------------|
| `GET`    | `/healthz`                  | —                                                           | Health check                   |
| `GET`    | `/api/boards`               | —                                                           | List all boards                |
| `POST`   | `/api/boards`               | `{"name": "..."}`                                           | Create a board                 |
| `GET`    | `/api/boards/:id`           | —                                                           | Get board with columns & tasks |
| `POST`   | `/api/boards/:id/columns`   | `{"name": "...", "position": 0}`                            | Add column to board            |
| `POST`   | `/api/columns/:id/tasks`    | `{"title": "...", "description": "...", "position": 0}`     | Add task to column             |
| `PUT`    | `/api/tasks/:id`            | `{"title": "...", "description": "...", "column_id": 1, "position": 0}` | Update task       |
| `DELETE` | `/api/tasks/:id`            | —                                                           | Delete task                    |

## Testing

```bash
make test
```

All tests run against in-memory SQLite — no external dependencies, no cleanup needed.

## Configuration

| Env Var   | Default        | Description              |
|-----------|----------------|--------------------------|
| `PORT`    | `8080`         | HTTP listen port         |
| `DB_PATH` | `taskboard.db` | Path to SQLite database  |

## Build

```bash
make build
./taskboard-api
```
