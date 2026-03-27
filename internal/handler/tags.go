package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/digishoe/taskboard-api/internal/model"
	"github.com/digishoe/taskboard-api/internal/store"
)

// TagHandler handles tag-related HTTP requests.
type TagHandler struct {
	store *store.SQLiteStore
}

// NewTagHandler creates a new TagHandler.
func NewTagHandler(s *store.SQLiteStore) *TagHandler {
	return &TagHandler{store: s}
}

// List handles GET /api/tags.
func (h *TagHandler) List(w http.ResponseWriter, r *http.Request) {
	tags, err := h.store.ListTags()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if tags == nil {
		tags = []model.Tag{}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tags)
}

// Create handles POST /api/tags.
func (h *TagHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name  string `json:"name"`
		Color string `json:"color"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	if req.Color == "" {
		req.Color = "#6366f1"
	}
	tag, err := h.store.CreateTag(req.Name, req.Color)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(tag)
}

// Update handles PUT /api/tags/{id}.
func (h *TagHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid tag id", http.StatusBadRequest)
		return
	}
	var req struct {
		Name  string `json:"name"`
		Color string `json:"color"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	tag, err := h.store.UpdateTag(id, req.Name, req.Color)
	if err != nil {
		http.Error(w, "tag not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tag)
}

// Delete handles DELETE /api/tags/{id}.
func (h *TagHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid tag id", http.StatusBadRequest)
		return
	}
	if err := h.store.DeleteTag(id); err != nil {
		http.Error(w, "tag not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// SetTaskTags handles PUT /api/tasks/{id}/tags.
func (h *TagHandler) SetTaskTags(w http.ResponseWriter, r *http.Request) {
	taskID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid task id", http.StatusBadRequest)
		return
	}
	var req struct {
		TagIDs []int `json:"tag_ids"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	if req.TagIDs == nil {
		req.TagIDs = []int{}
	}
	if err := h.store.SetTaskTags(taskID, req.TagIDs); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
