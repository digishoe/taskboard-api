package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/digishoe/taskboard-api/internal/model"
)

func TestCreateColumn(t *testing.T) {
	s := newTestStore(t)
	h := NewColumnHandler(s)

	body := strings.NewReader(`{"name":"Blocked","position":3}`)
	req := httptest.NewRequest(http.MethodPost, "/api/boards/1/columns", body)
	req.SetPathValue("id", "1")
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.Create(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}
	var col model.Column
	json.NewDecoder(w.Body).Decode(&col)
	if col.Name != "Blocked" {
		t.Errorf("expected 'Blocked', got %q", col.Name)
	}
	if col.BoardID != 1 {
		t.Errorf("expected board_id 1, got %d", col.BoardID)
	}
}

func TestCreateColumnBadJSON(t *testing.T) {
	s := newTestStore(t)
	h := NewColumnHandler(s)

	body := strings.NewReader(`not json`)
	req := httptest.NewRequest(http.MethodPost, "/api/boards/1/columns", body)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()
	h.Create(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestCreateColumnBadBoardID(t *testing.T) {
	s := newTestStore(t)
	h := NewColumnHandler(s)

	body := strings.NewReader(`{"name":"X","position":0}`)
	req := httptest.NewRequest(http.MethodPost, "/api/boards/abc/columns", body)
	req.SetPathValue("id", "abc")
	w := httptest.NewRecorder()
	h.Create(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}
