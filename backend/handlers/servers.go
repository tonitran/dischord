package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/tonitran/dischord/models"
	"github.com/tonitran/dischord/store"
)

type ServerHandler struct {
	Store *store.Store
}

func (h *ServerHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name    string `json:"name"`
		OwnerID string `json:"owner_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if req.Name == "" || req.OwnerID == "" {
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
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
	writeJSON(w, http.StatusCreated, srv)
}

func (h *ServerHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	srv, err := h.Store.GetServer(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, srv)
}

func (h *ServerHandler) Join(w http.ResponseWriter, r *http.Request) {
	serverID := r.PathValue("id")

	var req struct {
		UserID string `json:"user_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if req.UserID == "" {
		http.Error(w, "user_id is required", http.StatusBadRequest)
		return
	}
	if err := h.Store.JoinServer(serverID, req.UserID); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "joined"})
}

func (h *ServerHandler) ListMembers(w http.ResponseWriter, r *http.Request) {
	serverID := r.PathValue("id")
	members, err := h.Store.GetServerMembers(serverID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, members)
}
