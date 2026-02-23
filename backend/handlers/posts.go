package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/tonitran/dischord/models"
	"github.com/tonitran/dischord/store"
)

type PostHandler struct {
	Store *store.Store
}

func (h *PostHandler) Create(w http.ResponseWriter, r *http.Request) {
	server_id := r.PathValue("server_id")
	var req struct {
		AuthorID string `json:"author_id"`
		Title    string `json:"title"`
		Body     string `json:"body"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("posts: Create: failed to decode request body", "server_id", server_id, "error", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if req.AuthorID == "" || req.Title == "" || req.Body == "" {
		logger.Warn("posts: Create: missing required fields", "server_id", server_id, "author_id", req.AuthorID, "title", req.Title)
		http.Error(w, "author_id, title, and body are required", http.StatusBadRequest)
		return
	}

	now := time.Now()
	post := models.Post{
		ID:        generateID(),
		ServerID:  server_id,
		AuthorID:  req.AuthorID,
		Title:     req.Title,
		Body:      req.Body,
		Votes:     0,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := h.Store.CreatePost(post); err != nil {
		logger.Error("posts: Create: store error", "server_id", server_id, "author_id", req.AuthorID, "error", err)
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
	logger.Info("posts: Create: post created", "id", post.ID, "server_id", server_id, "author_id", req.AuthorID, "title", req.Title)
	writeJSON(w, http.StatusCreated, post)
}

func (h *PostHandler) Get(w http.ResponseWriter, r *http.Request) {
	server_id := r.PathValue("server_id")
	id := r.PathValue("id")
	logger.Debug("posts: Get: request", "server_id", server_id, "id", id)
	post, err := h.Store.GetPost(server_id, id)
	if err != nil {
		logger.Error("posts: Get: not found", "server_id", server_id, "id", id, "error", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	logger.Debug("posts: Get: success", "id", id, "title", post.Title)
	writeJSON(w, http.StatusOK, post)
}

func (h *PostHandler) Update(w http.ResponseWriter, r *http.Request) {
	server_id := r.PathValue("server_id")
	id := r.PathValue("id")
	post, err := h.Store.GetPost(server_id, id)
	if err != nil {
		logger.Error("posts: Update: post not found", "server_id", server_id, "id", id, "error", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	var req struct {
		Title *string `json:"title"`
		Body  *string `json:"body"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("posts: Update: failed to decode request body", "id", id, "error", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if req.Title != nil {
		post.Title = *req.Title
	}
	if req.Body != nil {
		post.Body = *req.Body
	}
	post.UpdatedAt = time.Now()

	if err := h.Store.UpdatePost(post); err != nil {
		logger.Error("posts: Update: store error", "id", id, "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	logger.Info("posts: Update: post updated", "id", id, "title", post.Title)
	writeJSON(w, http.StatusOK, post)
}

func (h *PostHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	logger.Debug("posts: Delete: request", "id", id)
	if err := h.Store.DeletePost(id); err != nil {
		logger.Error("posts: Delete: store error", "id", id, "error", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	logger.Info("posts: Delete: post deleted", "id", id)
	w.WriteHeader(http.StatusNoContent)
}
