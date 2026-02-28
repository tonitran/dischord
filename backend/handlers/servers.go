package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/tonitran/dischord/models"
	"github.com/tonitran/dischord/store"
)

type ServerHandler struct {
	Store *store.Database
}

func (h *ServerHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name    string `json:"name"`
		OwnerID string `json:"owner_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("servers: Create: failed to decode request body", "error", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if req.Name == "" || req.OwnerID == "" {
		logger.Warn("servers: Create: missing required fields", "name", req.Name, "owner_id", req.OwnerID)
		http.Error(w, "name and owner_id are required", http.StatusBadRequest)
		return
	}

	srv := models.Server{
		ID:        generateID(),
		Name:      req.Name,
		OwnerID:   req.OwnerID,
		MemberIDs: []string{req.OwnerID},
		CreatedAt: time.Now(),
	}
	if err := h.Store.CreateServer(srv); err != nil {
		logger.Error("servers: Create: store error", "name", req.Name, "owner_id", req.OwnerID, "error", err)
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
	if err := h.Store.JoinServer(srv.ID, srv.OwnerID); err != nil {
		logger.Error("servers: Create: failed to auto-join owner", "server_id", srv.ID, "owner_id", srv.OwnerID, "error", err)
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
	logger.Info("servers: Create: server created", "id", srv.ID, "name", srv.Name, "owner_id", srv.OwnerID)
	writeJSON(w, http.StatusCreated, srv)
}

func (h *ServerHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	logger.Debug("servers: Get: request", "id", id)
	srv, err := h.Store.GetServer(id)
	if err != nil {
		logger.Error("servers: Get: not found", "id", id, "error", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	logger.Debug("servers: Get: success", "id", id, "name", srv.Name)
	writeJSON(w, http.StatusOK, srv)
}

func (h *ServerHandler) Join(w http.ResponseWriter, r *http.Request) {
	serverID := r.PathValue("id")

	var req struct {
		UserID string `json:"user_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("servers: Join: failed to decode request body", "server_id", serverID, "error", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if req.UserID == "" {
		logger.Warn("servers: Join: missing user_id", "server_id", serverID)
		http.Error(w, "user_id is required", http.StatusBadRequest)
		return
	}
	if err := h.Store.JoinServer(serverID, req.UserID); err != nil {
		logger.Error("servers: Join: store error", "server_id", serverID, "user_id", req.UserID, "error", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	logger.Info("servers: Join: user joined server", "server_id", serverID, "user_id", req.UserID)
	writeJSON(w, http.StatusOK, map[string]string{"status": "joined"})
}

func (h *ServerHandler) ListMembers(w http.ResponseWriter, r *http.Request) {
	serverID := r.PathValue("id")
	logger.Debug("servers: ListMembers: request", "server_id", serverID)
	members, err := h.Store.GetServerMembers(serverID)
	if err != nil {
		logger.Error("servers: ListMembers: store error", "server_id", serverID, "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	logger.Debug("servers: ListMembers: success", "server_id", serverID, "count", len(members))
	writeJSON(w, http.StatusOK, members)
}
