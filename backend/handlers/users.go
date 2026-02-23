package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"time"

	"github.com/tonitran/dischord/models"
	"github.com/tonitran/dischord/store"
)

type UserHandler struct {
	Store *store.Store
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("users: Create: failed to decode request body", "error", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if req.Username == "" || req.Email == "" {
		logger.Warn("users: Create: missing required fields", "username", req.Username, "email", req.Email)
		http.Error(w, "username and email are required", http.StatusBadRequest)
		return
	}

	user := models.User{
		ID:        generateID(),
		Username:  req.Username,
		Email:     req.Email,
		CreatedAt: time.Now(),
	}
	if err := h.Store.CreateUser(user); err != nil {
		logger.Error("users: Create: store error", "username", req.Username, "email", req.Email, "error", err)
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
	logger.Info("users: Create: user created", "id", user.ID, "username", user.Username)
	writeJSON(w, http.StatusCreated, user)
}

func (h *UserHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	logger.Debug("users: Get: request", "id", id)
	user, err := h.Store.GetUser(id)
	if err != nil {
		logger.Error("users: Get: not found", "id", id, "error", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	logger.Debug("users: Get: success", "id", id, "username", user.Username)
	writeJSON(w, http.StatusOK, user)
}

func generateID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
