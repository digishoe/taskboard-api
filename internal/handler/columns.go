package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/digishoe/taskboard-api/internal/store"
)

// ColumnHandler handles column-related HTTP requests.
type ColumnHandler struct {
	store *store.SQLiteStore
}

// NewColumnHandler creates a new ColumnHandler.
func NewColumnHandler(s *store.SQLiteStore) *ColumnHandler {
	return &ColumnHandler{store: s}
}

// Create handles POST /api/boards/{id}/columns.
func (h *ColumnHandler) Create(w http.ResponseWriter, r *http.Request) {
	boardID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid board id", http.StatusBadRequest)
		return
	}
	var req struct {
		Name     string `json:"name"`
		Position int    `json:"position"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	col, err := h.store.CreateColumn(boardID, req.Name, req.Position)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(col)
}
