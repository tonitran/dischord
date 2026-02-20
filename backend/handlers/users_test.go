package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/tonitran/dischord/models"
)

func TestUserHandler_Create(t *testing.T) {
	h := &UserHandler{Store: testStore(t)}

	tests := []struct {
		name       string
		body       string
		wantStatus int
	}{
		{
			name:       "valid user",
			body:       `{"username":"alice","email":"alice@example.com"}`,
			wantStatus: http.StatusCreated,
		},
		{
			name:       "missing username",
			body:       `{"email":"alice@example.com"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "missing email",
			body:       `{"username":"alice"}`,
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
			req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(tt.body))
			w := httptest.NewRecorder()

			h.Create(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("got status %d, want %d", w.Code, tt.wantStatus)
			}
			if tt.wantStatus == http.StatusCreated {
				var user models.User
				if err := json.NewDecoder(w.Body).Decode(&user); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				if user.ID == "" {
					t.Error("expected non-empty user ID")
				}
				if user.Username != "alice" {
					t.Errorf("got username %q, want %q", user.Username, "alice")
				}
			}
		})
	}
}

func TestUserHandler_Get(t *testing.T) {
	s := testStore(t)
	h := &UserHandler{Store: s}

	user := models.User{ID: "u1", Username: "alice", Email: "alice@example.com"}
	s.CreateUser(user)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /users/{id}", h.Get)

	t.Run("existing user", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/users/u1", nil)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("got status %d, want %d", w.Code, http.StatusOK)
		}
		var got models.User
		json.NewDecoder(w.Body).Decode(&got)
		if got.Username != "alice" {
			t.Errorf("got username %q, want %q", got.Username, "alice")
		}
	})

	t.Run("nonexistent user", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/users/missing", nil)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("got status %d, want %d", w.Code, http.StatusNotFound)
		}
	})
}
