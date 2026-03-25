package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/digishoe/taskboard-api/internal/model"
	"github.com/digishoe/taskboard-api/internal/store"
)

func newTestStore(t *testing.T) *store.SQLiteStore {
	t.Helper()
	s, err := store.New(":memory:")
	if err != nil {
		t.Fatalf("store: %v", err)
	}
	t.Cleanup(func() { s.Close() })
	return s
}

func TestListBoards(t *testing.T) {
	s := newTestStore(t)
	h := NewBoardHandler(s)

	req := httptest.NewRequest(http.MethodGet, "/api/boards", nil)
	w := httptest.NewRecorder()
	h.List(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var boards []model.Board
	json.NewDecoder(w.Body).Decode(&boards)
	if len(boards) < 1 {
		t.Error("expected at least 1 board from seed")
	}
}

func TestCreateBoard(t *testing.T) {
	s := newTestStore(t)
	h := NewBoardHandler(s)

	body := strings.NewReader(`{"name":"Sprint 1"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/boards", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.Create(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}
	var b model.Board
	json.NewDecoder(w.Body).Decode(&b)
	if b.Name != "Sprint 1" {
		t.Errorf("expected 'Sprint 1', got %q", b.Name)
	}
}

func TestCreateBoardBadJSON(t *testing.T) {
	s := newTestStore(t)
	h := NewBoardHandler(s)

	body := strings.NewReader(`{bad`)
	req := httptest.NewRequest(http.MethodPost, "/api/boards", body)
	w := httptest.NewRecorder()
	h.Create(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestGetBoard(t *testing.T) {
	s := newTestStore(t)
	h := NewBoardHandler(s)

	req := httptest.NewRequest(http.MethodGet, "/api/boards/1", nil)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()
	h.Get(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var b model.Board
	json.NewDecoder(w.Body).Decode(&b)
	if len(b.Columns) != 3 {
		t.Errorf("expected 3 columns, got %d", len(b.Columns))
	}
}

func TestGetBoardNotFound(t *testing.T) {
	s := newTestStore(t)
	h := NewBoardHandler(s)

	req := httptest.NewRequest(http.MethodGet, "/api/boards/999", nil)
	req.SetPathValue("id", "999")
	w := httptest.NewRecorder()
	h.Get(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}
