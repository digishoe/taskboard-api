package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/digishoe/taskboard-api/internal/store"
)

// BoardHandler handles board-related HTTP requests.
type BoardHandler struct {
	store *store.SQLiteStore
}

// NewBoardHandler creates a new BoardHandler.
func NewBoardHandler(s *store.SQLiteStore) *BoardHandler {
	return &BoardHandler{store: s}
}

// List handles GET /api/boards.
func (h *BoardHandler) List(w http.ResponseWriter, r *http.Request) {
	boards, err := h.store.ListBoards()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(boards)
}

// Create handles POST /api/boards.
func (h *BoardHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	board, err := h.store.CreateBoard(req.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(board)
}

// Get handles GET /api/boards/{id}.
func (h *BoardHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid board id", http.StatusBadRequest)
		return
	}
	board, err := h.store.GetBoard(id)
	if err != nil {
		http.Error(w, "board not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(board)
}
