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
