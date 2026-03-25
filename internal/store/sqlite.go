package store

import (
	"database/sql"
	"fmt"

	"github.com/digishoe/taskboard-api/internal/model"

	_ "modernc.org/sqlite"
)

// SQLiteStore manages all database operations.
type SQLiteStore struct {
	db *sql.DB
}

// New opens (or creates) a SQLite database and runs migrations and seed data.
// Use ":memory:" for an ephemeral in-memory database (useful for tests).
func New(dsn string) (*SQLiteStore, error) {
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}
	// Enable WAL mode and foreign keys.
	for _, pragma := range []string{
		"PRAGMA journal_mode=WAL",
		"PRAGMA foreign_keys=ON",
	} {
		if _, err := db.Exec(pragma); err != nil {
			db.Close()
			return nil, fmt.Errorf("exec %q: %w", pragma, err)
		}
	}
	s := &SQLiteStore{db: db}
	if err := s.migrate(); err != nil {
		db.Close()
		return nil, fmt.Errorf("migrate: %w", err)
	}
	if err := s.seed(); err != nil {
		db.Close()
		return nil, fmt.Errorf("seed: %w", err)
	}
	return s, nil
}

// Close closes the underlying database connection.
func (s *SQLiteStore) Close() error {
	return s.db.Close()
}

func (s *SQLiteStore) migrate() error {
	schema := `
	CREATE TABLE IF NOT EXISTS boards (
		id   INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL
	);
	CREATE TABLE IF NOT EXISTS columns (
		id       INTEGER PRIMARY KEY AUTOINCREMENT,
		board_id INTEGER NOT NULL REFERENCES boards(id),
		name     TEXT NOT NULL,
		position INTEGER NOT NULL DEFAULT 0
	);
	CREATE TABLE IF NOT EXISTS tasks (
		id          INTEGER PRIMARY KEY AUTOINCREMENT,
		column_id   INTEGER NOT NULL REFERENCES columns(id),
		title       TEXT NOT NULL,
		description TEXT NOT NULL DEFAULT '',
		position    INTEGER NOT NULL DEFAULT 0
	);`
	_, err := s.db.Exec(schema)
	return err
}

func (s *SQLiteStore) seed() error {
	var count int
	if err := s.db.QueryRow("SELECT COUNT(*) FROM boards").Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return nil
	}
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	res, err := tx.Exec("INSERT INTO boards (name) VALUES (?)", "Default Board")
	if err != nil {
		return err
	}
	boardID, _ := res.LastInsertId()

	for i, name := range []string{"To Do", "In Progress", "Done"} {
		if _, err := tx.Exec("INSERT INTO columns (board_id, name, position) VALUES (?, ?, ?)", boardID, name, i); err != nil {
			return err
		}
	}
	return tx.Commit()
}

// --- Board operations ---

// ListBoards returns all boards (without nested columns/tasks).
func (s *SQLiteStore) ListBoards() ([]model.Board, error) {
	rows, err := s.db.Query("SELECT id, name FROM boards ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var boards []model.Board
	for rows.Next() {
		var b model.Board
		if err := rows.Scan(&b.ID, &b.Name); err != nil {
			return nil, err
		}
		boards = append(boards, b)
	}
	return boards, rows.Err()
}

// CreateBoard inserts a new board and returns it.
func (s *SQLiteStore) CreateBoard(name string) (model.Board, error) {
	res, err := s.db.Exec("INSERT INTO boards (name) VALUES (?)", name)
	if err != nil {
		return model.Board{}, err
	}
	id, _ := res.LastInsertId()
	return model.Board{ID: int(id), Name: name}, nil
}

// GetBoard returns a board with its columns and tasks nested.
func (s *SQLiteStore) GetBoard(id int) (model.Board, error) {
	var b model.Board
	err := s.db.QueryRow("SELECT id, name FROM boards WHERE id = ?", id).Scan(&b.ID, &b.Name)
	if err != nil {
		return b, fmt.Errorf("board %d not found: %w", id, err)
	}

	colRows, err := s.db.Query("SELECT id, board_id, name, position FROM columns WHERE board_id = ? ORDER BY position", id)
	if err != nil {
		return b, err
	}
	defer colRows.Close()

	for colRows.Next() {
		var c model.Column
		if err := colRows.Scan(&c.ID, &c.BoardID, &c.Name, &c.Position); err != nil {
			return b, err
		}
		b.Columns = append(b.Columns, c)
	}
	if err := colRows.Err(); err != nil {
		return b, err
	}

	for i := range b.Columns {
		taskRows, err := s.db.Query("SELECT id, column_id, title, description, position FROM tasks WHERE column_id = ? ORDER BY position", b.Columns[i].ID)
		if err != nil {
			return b, err
		}
		for taskRows.Next() {
			var tk model.Task
			if err := taskRows.Scan(&tk.ID, &tk.ColumnID, &tk.Title, &tk.Description, &tk.Position); err != nil {
				taskRows.Close()
				return b, err
			}
			b.Columns[i].Tasks = append(b.Columns[i].Tasks, tk)
		}
		taskRows.Close()
		if err := taskRows.Err(); err != nil {
			return b, err
		}
	}
	return b, nil
}

// --- Column operations ---

// CreateColumn inserts a new column into a board.
func (s *SQLiteStore) CreateColumn(boardID int, name string, position int) (model.Column, error) {
	res, err := s.db.Exec("INSERT INTO columns (board_id, name, position) VALUES (?, ?, ?)", boardID, name, position)
	if err != nil {
		return model.Column{}, err
	}
	id, _ := res.LastInsertId()
	return model.Column{ID: int(id), BoardID: boardID, Name: name, Position: position}, nil
}

// --- Task operations ---

// CreateTask inserts a new task into a column.
func (s *SQLiteStore) CreateTask(columnID int, title, description string, position int) (model.Task, error) {
	res, err := s.db.Exec("INSERT INTO tasks (column_id, title, description, position) VALUES (?, ?, ?, ?)", columnID, title, description, position)
	if err != nil {
		return model.Task{}, err
	}
	id, _ := res.LastInsertId()
	return model.Task{ID: int(id), ColumnID: columnID, Title: title, Description: description, Position: position}, nil
}

// UpdateTask updates an existing task's fields.
func (s *SQLiteStore) UpdateTask(id int, t model.Task) (model.Task, error) {
	res, err := s.db.Exec(
		"UPDATE tasks SET title = ?, description = ?, column_id = ?, position = ? WHERE id = ?",
		t.Title, t.Description, t.ColumnID, t.Position, id,
	)
	if err != nil {
		return model.Task{}, err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return model.Task{}, fmt.Errorf("task %d not found", id)
	}
	t.ID = id
	return t, nil
}

// DeleteTask removes a task by ID.
func (s *SQLiteStore) DeleteTask(id int) error {
	res, err := s.db.Exec("DELETE FROM tasks WHERE id = ?", id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("task %d not found", id)
	}
	return nil
}
