package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/digishoe/taskboard-api/internal/model"
	"github.com/digishoe/taskboard-api/internal/store"
)

// TaskHandler handles task-related HTTP requests.
type TaskHandler struct {
	store *store.SQLiteStore
}

// NewTaskHandler creates a new TaskHandler.
func NewTaskHandler(s *store.SQLiteStore) *TaskHandler {
	return &TaskHandler{store: s}
}

// Create handles POST /api/columns/{id}/tasks.
func (h *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	columnID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid column id", http.StatusBadRequest)
		return
	}
	var req struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Position    int    `json:"position"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	task, err := h.store.CreateTask(columnID, req.Title, req.Description, req.Position)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

// Update handles PUT /api/tasks/{id}.
func (h *TaskHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid task id", http.StatusBadRequest)
		return
	}
	var req model.Task
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	task, err := h.store.UpdateTask(id, req)
	if err != nil {
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

// Delete handles DELETE /api/tasks/{id}.
func (h *TaskHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid task id", http.StatusBadRequest)
		return
	}
	if err := h.store.DeleteTask(id); err != nil {
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
