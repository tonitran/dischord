package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/tonitran/dischord/models"
	"github.com/tonitran/dischord/store"
)

type MessageHandler struct {
	Store *store.Store
}

func (h *MessageHandler) Create(w http.ResponseWriter, r *http.Request) {
	serverID := r.PathValue("server_id")

	var req struct {
		AuthorID string `json:"author_id"`
		Content  string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if req.AuthorID == "" || req.Content == "" {
		http.Error(w, "author_id and content are required", http.StatusBadRequest)
		return
	}

	msg := models.Message{
		ID:        generateID(),
		ServerID:  serverID,
		AuthorID:  req.AuthorID,
		Content:   req.Content,
		CreatedAt: time.Now(),
	}
	if err := h.Store.CreateMessage(msg); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusCreated, msg)
}

func (h *MessageHandler) ListByServer(w http.ResponseWriter, r *http.Request) {
	serverID := r.PathValue("server_id")
	msgs := h.Store.GetMessagesByServer(serverID)
	writeJSON(w, http.StatusOK, msgs)
}
