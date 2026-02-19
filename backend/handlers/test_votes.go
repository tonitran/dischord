package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/tonitran/dischord/models"
	"github.com/tonitran/dischord/store"
)

func setupVotesTest() (*store.Store, *http.ServeMux) {
	s := store.New()
	h := &VoteHandler{Store: s}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /servers/{server_id}/posts/{id}/vote", h.GetVote)
	mux.HandleFunc("PUT /servers/{server_id}/posts/{id}/vote", h.PutVote)
	return s, mux
}

func TestPostHandler_GetVote(t *testing.T) {
	s, mux := setupVotesTest()
	s.CreatePost(models.Post{ID: "p1", ServerID: "s1", AuthorID: "u1", Title: "Hello", Body: "World"})
	s.PostVote("p1", "u1", 1)

	t.Run("existing vote", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/servers/s1/posts/p1/vote", strings.NewReader(`{"author":"u1"}`))
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("got status %d, want %d", w.Code, http.StatusOK)
		}
		var vote models.Vote
		json.NewDecoder(w.Body).Decode(&vote)
		if vote.Vote != 1 {
			t.Errorf("got vote %d, want %d", vote.Vote, 1)
		}
	})

	t.Run("nonexistent vote", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/servers/s1/posts/p1/vote", strings.NewReader(`{"author":"u2"}`))
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("got status %d, want %d", w.Code, http.StatusNotFound)
		}
	})

	t.Run("invalid json", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/servers/s1/posts/p1/vote", strings.NewReader(`{bad`))
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("got status %d, want %d", w.Code, http.StatusBadRequest)
		}
	})
}

func TestPostHandler_PutVote(t *testing.T) {
	s, mux := setupVotesTest()
	s.CreatePost(models.Post{ID: "p1", ServerID: "s1", AuthorID: "u1", Title: "Hello", Body: "World"})

	t.Run("upvote", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/servers/s1/posts/p1/vote", strings.NewReader(`{"author":"u1","vote":1}`))
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("got status %d, want %d", w.Code, http.StatusOK)
		}
		var post models.Post
		json.NewDecoder(w.Body).Decode(&post)
		if post.ID != "p1" {
			t.Errorf("got post ID %q, want %q", post.ID, "p1")
		}
	})

	t.Run("downvote", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/servers/s1/posts/p1/vote", strings.NewReader(`{"author":"u2","vote":-1}`))
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("got status %d, want %d", w.Code, http.StatusOK)
		}
	})

	t.Run("same vote is no-op", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/servers/s1/posts/p1/vote", strings.NewReader(`{"author":"u1","vote":1}`))
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("got status %d, want %d", w.Code, http.StatusOK)
		}
	})

	t.Run("nonexistent post", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/servers/s1/posts/missing/vote", strings.NewReader(`{"author":"u1","vote":1}`))
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("got status %d, want %d", w.Code, http.StatusNotFound)
		}
	})

	t.Run("invalid json", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/servers/s1/posts/p1/vote", strings.NewReader(`{bad`))
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("got status %d, want %d", w.Code, http.StatusBadRequest)
		}
	})
}
