package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/tonitran/dischord/store"
)

type VoteHandler struct {
	Store *store.Database
}

func (h *VoteHandler) GetVote(w http.ResponseWriter, r *http.Request) {
	post_id := r.PathValue("id")

	var req struct {
		Author string `json:"author"`
		Vote   int    `json:"vote"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("votes: GetVote: failed to decode request body", "post_id", post_id, "error", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	logger.Debug("votes: GetVote: request", "post_id", post_id, "author", req.Author)
	vote, err := h.Store.GetVote(post_id, req.Author)
	if err != nil {
		logger.Error("votes: GetVote: not found", "post_id", post_id, "author", req.Author, "error", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	logger.Debug("votes: GetVote: success", "post_id", post_id, "author", req.Author, "vote", vote.Vote)
	writeJSON(w, http.StatusOK, vote)
}

func (h *VoteHandler) PutVote(w http.ResponseWriter, r *http.Request) {
	server_id := r.PathValue("server_id")
	post_id := r.PathValue("id")
	post, err := h.Store.GetPost(server_id, post_id)
	if err != nil {
		logger.Error("votes: PutVote: post not found", "server_id", server_id, "post_id", post_id, "error", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	var req struct {
		Author string `json:"author"`
		Vote   int    `json:"vote"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("votes: PutVote: failed to decode request body", "post_id", post_id, "error", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	vote, _ := h.Store.GetVote(post_id, req.Author)
	if req.Author != "" && req.Vote >= -1 && req.Vote <= 1 && req.Vote != vote.Vote {
		if err := h.Store.PostVote(post_id, req.Author, req.Vote); err != nil {
			logger.Error("votes: PutVote: store error", "post_id", post_id, "author", req.Author, "vote", req.Vote, "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		logger.Info("votes: PutVote: vote recorded", "post_id", post_id, "author", req.Author, "vote", req.Vote)
	} else {
		logger.Debug("votes: PutVote: vote unchanged or invalid", "post_id", post_id, "author", req.Author, "requested_vote", req.Vote, "existing_vote", vote.Vote)
	}
	writeJSON(w, http.StatusOK, post)
}
