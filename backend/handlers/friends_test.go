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

func setupFriendsTest(t *testing.T) (*store.Database, *http.ServeMux) {
	s := testStore(t)
	h := &FriendHandler{Store: s}

	s.CreateUser(models.User{ID: "u1", Username: "alice", Email: "a@example.com"})
	s.CreateUser(models.User{ID: "u2", Username: "bob", Email: "b@example.com"})

	mux := http.NewServeMux()
	mux.HandleFunc("POST /users/{id}/friends", h.Add)
	mux.HandleFunc("GET /users/{id}/friends", h.List)
	return s, mux
}

func TestFriendHandler_Add(t *testing.T) {
	_, mux := setupFriendsTest(t)

	tests := []struct {
		name       string
		userID     string
		body       string
		wantStatus int
	}{
		{
			name:       "valid add",
			userID:     "u1",
			body:       `{"friend_id":"u2"}`,
			wantStatus: http.StatusOK,
		},
		{
			name:       "missing friend_id",
			userID:     "u1",
			body:       `{}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "nonexistent user",
			userID:     "u1",
			body:       `{"friend_id":"missing"}`,
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "invalid json",
			userID:     "u1",
			body:       `{bad`,
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/users/"+tt.userID+"/friends", strings.NewReader(tt.body))
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("got status %d, want %d", w.Code, tt.wantStatus)
			}
		})
	}
}

func TestFriendHandler_List(t *testing.T) {
	s, mux := setupFriendsTest(t)
	s.AddFriend("u1", "u2")

	t.Run("user with friends", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/users/u1/friends", nil)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("got status %d, want %d", w.Code, http.StatusOK)
		}
		var friends []models.User
		json.NewDecoder(w.Body).Decode(&friends)
		if len(friends) != 1 {
			t.Errorf("got %d friends, want 1", len(friends))
		}
	})

	t.Run("user with no friends", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/users/nobody/friends", nil)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("got status %d, want %d", w.Code, http.StatusOK)
		}
	})
}
