package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/tonitran/dischord/store"
)

type FriendHandler struct {
	Store *store.Store
}

func (h *FriendHandler) Add(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("id")

	var req struct {
		FriendID string `json:"friend_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if req.FriendID == "" {
		http.Error(w, "friend_id is required", http.StatusBadRequest)
		return
	}
	if err := h.Store.AddFriend(userID, req.FriendID); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "friend added"})
}

func (h *FriendHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("id")
	friends, err := h.Store.GetFriends(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, friends)
}
