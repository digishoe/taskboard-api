package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/digishoe/taskboard-api/internal/model"
)

func TestCreateTask(t *testing.T) {
	s := newTestStore(t)
	h := NewTaskHandler(s)

	body := strings.NewReader(`{"title":"Fix bug","description":"It crashes","position":0}`)
	req := httptest.NewRequest(http.MethodPost, "/api/columns/1/tasks", body)
	req.SetPathValue("id", "1")
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.Create(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}
	var tk model.Task
	json.NewDecoder(w.Body).Decode(&tk)
	if tk.Title != "Fix bug" {
		t.Errorf("expected 'Fix bug', got %q", tk.Title)
	}
}

func TestUpdateTask(t *testing.T) {
	s := newTestStore(t)
	th := NewTaskHandler(s)

	// First create a task.
	task, _ := s.CreateTask(1, "Old", "old desc", 0)

	body := strings.NewReader(`{"title":"New","description":"new desc","column_id":2,"position":1}`)
	req := httptest.NewRequest(http.MethodPut, "/api/tasks/1", body)
	req.SetPathValue("id", "1")
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	th.Update(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d; task id was %d", w.Code, task.ID)
	}
	var updated model.Task
	json.NewDecoder(w.Body).Decode(&updated)
	if updated.Title != "New" {
		t.Errorf("expected 'New', got %q", updated.Title)
	}
	if updated.ColumnID != 2 {
		t.Errorf("expected column_id 2, got %d", updated.ColumnID)
	}
}

func TestUpdateTaskNotFound(t *testing.T) {
	s := newTestStore(t)
	h := NewTaskHandler(s)

	body := strings.NewReader(`{"title":"X","description":"","column_id":1,"position":0}`)
	req := httptest.NewRequest(http.MethodPut, "/api/tasks/999", body)
	req.SetPathValue("id", "999")
	w := httptest.NewRecorder()
	h.Update(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestDeleteTask(t *testing.T) {
	s := newTestStore(t)
	h := NewTaskHandler(s)

	s.CreateTask(1, "Delete me", "", 0)

	req := httptest.NewRequest(http.MethodDelete, "/api/tasks/1", nil)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()
	h.Delete(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", w.Code)
	}
}

func TestDeleteTaskNotFound(t *testing.T) {
	s := newTestStore(t)
	h := NewTaskHandler(s)

	req := httptest.NewRequest(http.MethodDelete, "/api/tasks/999", nil)
	req.SetPathValue("id", "999")
	w := httptest.NewRecorder()
	h.Delete(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestCreateTaskBadJSON(t *testing.T) {
	s := newTestStore(t)
	h := NewTaskHandler(s)

	body := strings.NewReader(`{bad}`)
	req := httptest.NewRequest(http.MethodPost, "/api/columns/1/tasks", body)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()
	h.Create(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}
