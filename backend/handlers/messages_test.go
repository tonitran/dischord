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

func setupMessagesTest(t *testing.T) (*store.Store, *http.ServeMux) {
	s := testStore(t)
	h := &MessageHandler{Store: s}

	s.CreateServer(models.Server{ID: "s1", Name: "test-server", OwnerID: "u1", MemberIDs: []string{"u1"}})

	mux := http.NewServeMux()
	mux.HandleFunc("POST /servers/{server_id}/messages", h.Create)
	mux.HandleFunc("GET /servers/{server_id}/messages", h.ListByServer)
	return s, mux
}

func TestMessageHandler_Create(t *testing.T) {
	_, mux := setupMessagesTest(t)

	tests := []struct {
		name       string
		serverID   string
		body       string
		wantStatus int
	}{
		{
			name:       "valid message",
			serverID:   "s1",
			body:       `{"author_id":"u1","content":"hello"}`,
			wantStatus: http.StatusCreated,
		},
		{
			name:       "missing content",
			serverID:   "s1",
			body:       `{"author_id":"u1"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "missing author_id",
			serverID:   "s1",
			body:       `{"content":"hello"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "nonexistent server",
			serverID:   "missing",
			body:       `{"author_id":"u1","content":"hello"}`,
			wantStatus: http.StatusNotFound,
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
			req := httptest.NewRequest(http.MethodPost, "/servers/"+tt.serverID+"/messages", strings.NewReader(tt.body))
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("got status %d, want %d", w.Code, tt.wantStatus)
			}
			if tt.wantStatus == http.StatusCreated {
				var msg models.Message
				json.NewDecoder(w.Body).Decode(&msg)
				if msg.ID == "" {
					t.Error("expected non-empty message ID")
				}
				if msg.Content != "hello" {
					t.Errorf("got content %q, want %q", msg.Content, "hello")
				}
				if msg.ServerID != tt.serverID {
					t.Errorf("got server_id %q, want %q", msg.ServerID, tt.serverID)
				}
			}
		})
	}
}

func TestMessageHandler_ListByServer(t *testing.T) {
	s, mux := setupMessagesTest(t)

	s.CreateMessage(models.Message{ID: "m1", ServerID: "s1", AuthorID: "u1", Content: "hello"})
	s.CreateMessage(models.Message{ID: "m2", ServerID: "s1", AuthorID: "u1", Content: "world"})

	t.Run("server with messages", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/servers/s1/messages", nil)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("got status %d, want %d", w.Code, http.StatusOK)
		}
		var msgs []models.Message
		json.NewDecoder(w.Body).Decode(&msgs)
		if len(msgs) != 2 {
			t.Errorf("got %d messages, want 2", len(msgs))
		}
	})

	t.Run("server with no messages", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/servers/empty/messages", nil)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("got status %d, want %d", w.Code, http.StatusOK)
		}
	})
}
