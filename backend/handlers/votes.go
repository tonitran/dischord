package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/tonitran/dischord/store"
)

type VoteHandler struct {
	Store *store.Store
}

func (h *VoteHandler) GetVote(w http.ResponseWriter, r *http.Request) {
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

func (h *VoteHandler) PutVote(w http.ResponseWriter, r *http.Request) {
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
