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

func setupPostsTest(t *testing.T) (*store.Database, *http.ServeMux) {
	s := testStore(t)
	h := &PostHandler{Store: s}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /posts", h.Create)
	mux.HandleFunc("GET /posts/{id}", h.Get)
	mux.HandleFunc("PATCH /posts/{id}", h.Update)
	mux.HandleFunc("DELETE /posts/{id}", h.Delete)
	return s, mux
}

func TestPostHandler_Create(t *testing.T) {
	_, mux := setupPostsTest(t)

	tests := []struct {
		name       string
		body       string
		wantStatus int
	}{
		{
			name:       "valid post",
			body:       `{"author_id":"u1","title":"Hello","body":"World"}`,
			wantStatus: http.StatusCreated,
		},
		{
			name:       "missing title",
			body:       `{"author_id":"u1"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "missing author_id",
			body:       `{"title":"Hello"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid json",
			body:       `{bad`,
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/posts", strings.NewReader(tt.body))
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("got status %d, want %d", w.Code, tt.wantStatus)
			}
			if tt.wantStatus == http.StatusCreated {
				var post models.Post
				json.NewDecoder(w.Body).Decode(&post)
				if post.ID == "" {
					t.Error("expected non-empty post ID")
				}
				if post.Title != "Hello" {
					t.Errorf("got title %q, want %q", post.Title, "Hello")
				}
			}
		})
	}
}

func TestPostHandler_Get(t *testing.T) {
	s, mux := setupPostsTest(t)
	s.CreatePost(models.Post{ID: "p1", AuthorID: "u1", Title: "Hello", Body: "World"})

	t.Run("existing post", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/posts/p1", nil)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("got status %d, want %d", w.Code, http.StatusOK)
		}
		var post models.Post
		json.NewDecoder(w.Body).Decode(&post)
		if post.Title != "Hello" {
			t.Errorf("got title %q, want %q", post.Title, "Hello")
		}
		if post.Body != "World" {
			t.Errorf("got body %q, want %q", post.Body, "World")
		}
	})

	t.Run("nonexistent post", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/posts/missing", nil)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("got status %d, want %d", w.Code, http.StatusNotFound)
		}
	})
}

func TestPostHandler_Update(t *testing.T) {
	s, mux := setupPostsTest(t)
	s.CreatePost(models.Post{ID: "p1", AuthorID: "u1", Title: "Hello", Body: "World"})

	t.Run("update title", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPatch, "/posts/p1", strings.NewReader(`{"title":"Updated"}`))
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("got status %d, want %d", w.Code, http.StatusOK)
		}
		var post models.Post
		json.NewDecoder(w.Body).Decode(&post)
		if post.Title != "Updated" {
			t.Errorf("got title %q, want %q", post.Title, "Updated")
		}
		if post.Body != "World" {
			t.Errorf("got body %q, want %q", post.Body, "World")
		}
	})

	t.Run("update body", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPatch, "/posts/p1", strings.NewReader(`{"body":"New body"}`))
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("got status %d, want %d", w.Code, http.StatusOK)
		}
		var post models.Post
		json.NewDecoder(w.Body).Decode(&post)
		if post.Body != "New body" {
			t.Errorf("got body %q, want %q", post.Body, "New body")
		}
	})

	t.Run("nonexistent post", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPatch, "/posts/missing", strings.NewReader(`{"title":"X"}`))
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("got status %d, want %d", w.Code, http.StatusNotFound)
		}
	})

	t.Run("invalid json", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPatch, "/posts/p1", strings.NewReader(`{bad`))
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("got status %d, want %d", w.Code, http.StatusBadRequest)
		}
	})
}

func TestPostHandler_Delete(t *testing.T) {
	s, mux := setupPostsTest(t)
	s.CreatePost(models.Post{ID: "p1", AuthorID: "u1", Title: "Hello"})

	t.Run("existing post", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/posts/p1", nil)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		if w.Code != http.StatusNoContent {
			t.Errorf("got status %d, want %d", w.Code, http.StatusNoContent)
		}
	})

	t.Run("nonexistent post", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/posts/missing", nil)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("got status %d, want %d", w.Code, http.StatusNotFound)
		}
	})
}
