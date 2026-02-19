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
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if req.AuthorID == "" || req.Title == "" || req.Body == "" {
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
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
	writeJSON(w, http.StatusCreated, post)
}

func (h *PostHandler) Get(w http.ResponseWriter, r *http.Request) {
	server_id := r.PathValue("server_id")
	id := r.PathValue("id")
	post, err := h.Store.GetPost(server_id, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, post)
}

func (h *PostHandler) Update(w http.ResponseWriter, r *http.Request) {
	server_id := r.PathValue("server_id")
	id := r.PathValue("id")
	post, err := h.Store.GetPost(server_id, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	var req struct {
		Title *string `json:"title"`
		Body  *string `json:"body"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, post)
}

func (h *PostHandler) GetVote(w http.ResponseWriter, r *http.Request) {
	post_id := r.PathValue("id")

	var req struct {
		Author string `json:"author"`
		Vote   int    `json:"vote"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	vote, err := h.Store.GetVote(post_id, req.Author)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	writeJSON(w, http.StatusOK, vote)
}

func (h *PostHandler) PutVote(w http.ResponseWriter, r *http.Request) {
	server_id := r.PathValue("server_id")
	post_id := r.PathValue("id")
	post, err := h.Store.GetPost(server_id, post_id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	var req struct {
		Author string `json:"author"`
		Vote   int    `json:"vote"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	vote, _ := h.Store.GetVote(post_id, req.Author)
	if req.Author != "" && req.Vote >= -1 && req.Vote <= 1 && req.Vote != vote.Vote {
		if err := h.Store.PostVote(post_id, req.Author, req.Vote); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	writeJSON(w, http.StatusOK, post)
}

func (h *PostHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.Store.DeletePost(id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
