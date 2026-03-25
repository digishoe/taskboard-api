package store

import (
	"testing"

	"github.com/digishoe/taskboard-api/internal/model"
)

func mustNewTestStore(t *testing.T) *SQLiteStore {
	t.Helper()
	s, err := New(":memory:")
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}
	t.Cleanup(func() { s.Close() })
	return s
}

// --- Board tests ---

func TestSeedCreatesDefaultBoard(t *testing.T) {
	s := mustNewTestStore(t)
	boards, err := s.ListBoards()
	if err != nil {
		t.Fatalf("ListBoards: %v", err)
	}
	if len(boards) != 1 {
		t.Fatalf("expected 1 seeded board, got %d", len(boards))
	}
	if boards[0].Name != "Default Board" {
		t.Errorf("expected board name 'Default Board', got %q", boards[0].Name)
	}
}

func TestSeedCreatesThreeColumns(t *testing.T) {
	s := mustNewTestStore(t)
	board, err := s.GetBoard(1)
	if err != nil {
		t.Fatalf("GetBoard: %v", err)
	}
	if len(board.Columns) != 3 {
		t.Fatalf("expected 3 columns, got %d", len(board.Columns))
	}
	expected := []string{"To Do", "In Progress", "Done"}
	for i, col := range board.Columns {
		if col.Name != expected[i] {
			t.Errorf("column %d: expected %q, got %q", i, expected[i], col.Name)
		}
	}
}

func TestCreateBoard(t *testing.T) {
	s := mustNewTestStore(t)
	b, err := s.CreateBoard("My Board")
	if err != nil {
		t.Fatalf("CreateBoard: %v", err)
	}
	if b.ID == 0 {
		t.Error("expected non-zero ID")
	}
	if b.Name != "My Board" {
		t.Errorf("expected name 'My Board', got %q", b.Name)
	}
}

func TestListBoards(t *testing.T) {
	s := mustNewTestStore(t)
	s.CreateBoard("Board 2")
	boards, err := s.ListBoards()
	if err != nil {
		t.Fatalf("ListBoards: %v", err)
	}
	if len(boards) != 2 {
		t.Fatalf("expected 2 boards, got %d", len(boards))
	}
}

func TestGetBoardNotFound(t *testing.T) {
	s := mustNewTestStore(t)
	_, err := s.GetBoard(999)
	if err == nil {
		t.Fatal("expected error for non-existent board")
	}
}

// --- Column tests ---

func TestCreateColumn(t *testing.T) {
	s := mustNewTestStore(t)
	col, err := s.CreateColumn(1, "Blocked", 3)
	if err != nil {
		t.Fatalf("CreateColumn: %v", err)
	}
	if col.Name != "Blocked" {
		t.Errorf("expected 'Blocked', got %q", col.Name)
	}
	if col.BoardID != 1 {
		t.Errorf("expected board_id 1, got %d", col.BoardID)
	}
	if col.Position != 3 {
		t.Errorf("expected position 3, got %d", col.Position)
	}
}

// --- Task tests ---

func TestCreateTask(t *testing.T) {
	s := mustNewTestStore(t)
	task, err := s.CreateTask(1, "Fix bug", "Segfault on line 42", 0)
	if err != nil {
		t.Fatalf("CreateTask: %v", err)
	}
	if task.Title != "Fix bug" {
		t.Errorf("expected 'Fix bug', got %q", task.Title)
	}
	if task.ColumnID != 1 {
		t.Errorf("expected column_id 1, got %d", task.ColumnID)
	}
}

func TestUpdateTask(t *testing.T) {
	s := mustNewTestStore(t)
	task, _ := s.CreateTask(1, "Original", "Desc", 0)

	updated, err := s.UpdateTask(task.ID, model.Task{
		Title:       "Updated",
		Description: "New desc",
		ColumnID:    2,
		Position:    1,
	})
	if err != nil {
		t.Fatalf("UpdateTask: %v", err)
	}
	if updated.Title != "Updated" {
		t.Errorf("expected 'Updated', got %q", updated.Title)
	}
	if updated.ColumnID != 2 {
		t.Errorf("expected column_id 2, got %d", updated.ColumnID)
	}
}

func TestDeleteTask(t *testing.T) {
	s := mustNewTestStore(t)
	task, _ := s.CreateTask(1, "Delete me", "", 0)
	err := s.DeleteTask(task.ID)
	if err != nil {
		t.Fatalf("DeleteTask: %v", err)
	}
	// Verify board no longer contains this task
	board, _ := s.GetBoard(1)
	for _, col := range board.Columns {
		for _, tk := range col.Tasks {
			if tk.ID == task.ID {
				t.Error("task should have been deleted")
			}
		}
	}
}

func TestDeleteTaskNotFound(t *testing.T) {
	s := mustNewTestStore(t)
	err := s.DeleteTask(999)
	if err == nil {
		t.Error("expected error deleting non-existent task")
	}
}

func TestGetBoardIncludesTasks(t *testing.T) {
	s := mustNewTestStore(t)
	s.CreateTask(1, "Task A", "Desc A", 0)
	s.CreateTask(1, "Task B", "Desc B", 1)

	board, err := s.GetBoard(1)
	if err != nil {
		t.Fatalf("GetBoard: %v", err)
	}

	taskCount := 0
	for _, col := range board.Columns {
		taskCount += len(col.Tasks)
	}
	if taskCount != 2 {
		t.Errorf("expected 2 tasks, got %d", taskCount)
	}
}
