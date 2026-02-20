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

func setupServersTest(t *testing.T) (*store.Store, *http.ServeMux) {
	s := testStore(t)
	h := &ServerHandler{Store: s}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /servers", h.Create)
	mux.HandleFunc("GET /servers/{id}", h.Get)
	mux.HandleFunc("POST /servers/{id}/members", h.Join)
	mux.HandleFunc("GET /servers/{id}/members", h.ListMembers)
	return s, mux
}

func TestServerHandler_Create(t *testing.T) {
	_, mux := setupServersTest(t)

	tests := []struct {
		name       string
		body       string
		wantStatus int
	}{
		{
			name:       "valid server",
			body:       `{"name":"general","owner_id":"u1"}`,
			wantStatus: http.StatusCreated,
		},
		{
			name:       "missing name",
			body:       `{"owner_id":"u1"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "missing owner_id",
			body:       `{"name":"general"}`,
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
			req := httptest.NewRequest(http.MethodPost, "/servers", strings.NewReader(tt.body))
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("got status %d, want %d", w.Code, tt.wantStatus)
			}
			if tt.wantStatus == http.StatusCreated {
				var srv models.Server
				json.NewDecoder(w.Body).Decode(&srv)
				if srv.ID == "" {
					t.Error("expected non-empty server ID")
				}
				if srv.Name != "general" {
					t.Errorf("got name %q, want %q", srv.Name, "general")
				}
				if len(srv.MemberIDs) != 1 || srv.MemberIDs[0] != "u1" {
					t.Errorf("expected owner in member list, got %v", srv.MemberIDs)
				}
			}
		})
	}
}

func TestServerHandler_Join(t *testing.T) {
	s, mux := setupServersTest(t)
	s.CreateUser(models.User{ID: "u1", Username: "alice", Email: "a@example.com"})
	s.CreateServer(models.Server{ID: "s1", Name: "general", OwnerID: "u1"})

	tests := []struct {
		name       string
		serverID   string
		body       string
		wantStatus int
	}{
		{
			name:       "valid join",
			serverID:   "s1",
			body:       `{"user_id":"u1"}`,
			wantStatus: http.StatusOK,
		},
		{
			name:       "idempotent rejoin",
			serverID:   "s1",
			body:       `{"user_id":"u1"}`,
			wantStatus: http.StatusOK,
		},
		{
			name:       "nonexistent server",
			serverID:   "missing",
			body:       `{"user_id":"u1"}`,
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "nonexistent user",
			serverID:   "s1",
			body:       `{"user_id":"missing"}`,
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "missing user_id",
			serverID:   "s1",
			body:       `{}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid json",
			serverID:   "s1",
			body:       `{bad`,
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/servers/"+tt.serverID+"/members", strings.NewReader(tt.body))
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("got status %d, want %d", w.Code, tt.wantStatus)
			}
		})
	}
}

func TestServerHandler_ListMembers(t *testing.T) {
	s, mux := setupServersTest(t)
	s.CreateUser(models.User{ID: "u1", Username: "alice", Email: "a@example.com"})
	s.CreateUser(models.User{ID: "u2", Username: "bob", Email: "b@example.com"})
	s.CreateServer(models.Server{ID: "s1", Name: "general", OwnerID: "u1"})
	s.JoinServer("s1", "u1")
	s.JoinServer("s1", "u2")

	t.Run("lists all members", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/servers/s1/members", nil)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("got status %d, want %d", w.Code, http.StatusOK)
		}
		var members []models.User
		json.NewDecoder(w.Body).Decode(&members)
		if len(members) != 2 {
			t.Errorf("got %d members, want 2", len(members))
		}
	})

	t.Run("server with no members returns empty list", func(t *testing.T) {
		s.CreateServer(models.Server{ID: "s2", Name: "empty", OwnerID: "u1"})

		req := httptest.NewRequest(http.MethodGet, "/servers/s2/members", nil)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("got status %d, want %d", w.Code, http.StatusOK)
		}
		var members []models.User
		json.NewDecoder(w.Body).Decode(&members)
		if len(members) != 0 {
			t.Errorf("got %d members, want 0", len(members))
		}
	})
}

func TestServerHandler_Get(t *testing.T) {
	s, mux := setupServersTest(t)
	s.CreateServer(models.Server{ID: "s1", Name: "general", OwnerID: "u1", MemberIDs: []string{"u1"}})

	t.Run("existing server", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/servers/s1", nil)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("got status %d, want %d", w.Code, http.StatusOK)
		}
		var srv models.Server
		json.NewDecoder(w.Body).Decode(&srv)
		if srv.Name != "general" {
			t.Errorf("got name %q, want %q", srv.Name, "general")
		}
	})

	t.Run("nonexistent server", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/servers/missing", nil)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("got status %d, want %d", w.Code, http.StatusNotFound)
		}
	})
}
