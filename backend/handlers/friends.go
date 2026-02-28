package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/tonitran/dischord/store"
)

type FriendHandler struct {
	Store *store.Database
}

func (h *FriendHandler) Add(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("id")

	var req struct {
		FriendID string `json:"friend_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("friends: Add: failed to decode request body", "user_id", userID, "error", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if req.FriendID == "" {
		logger.Warn("friends: Add: missing friend_id", "user_id", userID)
		http.Error(w, "friend_id is required", http.StatusBadRequest)
		return
	}
	if err := h.Store.AddFriend(userID, req.FriendID); err != nil {
		logger.Error("friends: Add: store error", "user_id", userID, "friend_id", req.FriendID, "error", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	logger.Info("friends: Add: friend added", "user_id", userID, "friend_id", req.FriendID)
	writeJSON(w, http.StatusOK, map[string]string{"status": "friend added"})
}

func (h *FriendHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("id")
	logger.Debug("friends: List: request", "user_id", userID)
	friends, err := h.Store.GetFriends(userID)
	if err != nil {
		logger.Error("friends: List: store error", "user_id", userID, "error", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	logger.Debug("friends: List: success", "user_id", userID, "count", len(friends))
	writeJSON(w, http.StatusOK, friends)
}
