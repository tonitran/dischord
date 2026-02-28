package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/tonitran/dischord/models"
	"github.com/tonitran/dischord/store"
)

type MessageHandler struct {
	Store *store.Database
}

func (h *MessageHandler) Create(w http.ResponseWriter, r *http.Request) {
	serverID := r.PathValue("server_id")

	var req struct {
		AuthorID string `json:"author_id"`
		Content  string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("messages: Create: failed to decode request body", "server_id", serverID, "error", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if req.AuthorID == "" || req.Content == "" {
		logger.Warn("messages: Create: missing required fields", "server_id", serverID, "author_id", req.AuthorID)
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
		logger.Error("messages: Create: store error", "server_id", serverID, "author_id", req.AuthorID, "error", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	logger.Info("messages: Create: message created", "id", msg.ID, "server_id", serverID, "author_id", req.AuthorID)
	writeJSON(w, http.StatusCreated, msg)
}

func (h *MessageHandler) ListByServer(w http.ResponseWriter, r *http.Request) {
	serverID := r.PathValue("server_id")
	logger.Debug("messages: ListByServer: request", "server_id", serverID)
	msgs := h.Store.GetMessagesByServer(serverID)
	logger.Debug("messages: ListByServer: success", "server_id", serverID, "count", len(msgs))
	writeJSON(w, http.StatusOK, msgs)
}
