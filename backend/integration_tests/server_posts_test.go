package integration_tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/tonitran/dischord/models"
	"github.com/tonitran/dischord/router"
)

func TestServerPostIntegration(t *testing.T) {
	s := testStore(t)
	handler := router.New(s)

	// Step 0: Seed the owner user required by the FK constraint on server_user.
	if err := s.CreateUser(models.User{ID: "user-1", Username: "user1", Email: "user1@example.com"}); err != nil {
		t.Fatalf("seed user: %v", err)
	}

	// Step 1: Create a server.
	createServerBody := `{"name":"test-server","owner_id":"user-1"}`
	req := httptest.NewRequest(http.MethodPost, "/servers", strings.NewReader(createServerBody))
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("create server: got status %d, want %d\nbody: %s", w.Code, http.StatusCreated, w.Body.String())
	}
	var createdServer models.Server
	if err := json.NewDecoder(w.Body).Decode(&createdServer); err != nil {
		t.Fatalf("create server: failed to decode response: %v", err)
	}
	if createdServer.ID == "" {
		t.Fatal("create server: expected non-empty server ID")
	}
	t.Logf("created server with ID %q", createdServer.ID)

	// Step 2: Add a post to the server.
	createPostBody := `{"author_id":"user-1","title":"Hello World","body":"This is the first post."}`
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/servers/%s/posts", createdServer.ID), strings.NewReader(createPostBody))
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("create post: got status %d, want %d\nbody: %s", w.Code, http.StatusCreated, w.Body.String())
	}
	var createdPost models.Post
	if err := json.NewDecoder(w.Body).Decode(&createdPost); err != nil {
		t.Fatalf("create post: failed to decode response: %v", err)
	}
	if createdPost.ID == "" {
		t.Fatal("create post: expected non-empty post ID")
	}
	if createdPost.ServerID != createdServer.ID {
		t.Errorf("create post: got server_id %q, want %q", createdPost.ServerID, createdServer.ID)
	}
	t.Logf("created post with ID %q on server %q", createdPost.ID, createdPost.ServerID)

	// Step 3: Get the server and confirm the post is listed.
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/servers/%s", createdServer.ID), nil)
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("get server: got status %d, want %d\nbody: %s", w.Code, http.StatusOK, w.Body.String())
	}
	var fetchedServer models.Server
	if err := json.NewDecoder(w.Body).Decode(&fetchedServer); err != nil {
		t.Fatalf("get server: failed to decode response: %v", err)
	}

	found := false
	for _, postID := range fetchedServer.Posts {
		if postID == createdPost.ID {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("get server: post %q not found in server's posts list %v", createdPost.ID, fetchedServer.Posts)
	} else {
		t.Logf("confirmed post %q is listed in server %q posts", createdPost.ID, fetchedServer.ID)
	}

	// Step 4: Fetch the post directly and verify body.
	req = httptest.NewRequest(http.MethodGet,
		fmt.Sprintf("/servers/%s/posts/%s", createdServer.ID, createdPost.ID), nil)
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("get post: got status %d, want %d\nbody: %s", w.Code, http.StatusOK, w.Body.String())
	}
	var fetchedPost models.Post
	if err := json.NewDecoder(w.Body).Decode(&fetchedPost); err != nil {
		t.Fatalf("get post: failed to decode response: %v", err)
	}
	if fetchedPost.Body != "This is the first post." {
		t.Errorf("get post: got body %q, want %q", fetchedPost.Body, "This is the first post.")
	}
	t.Logf("confirmed post body %q", fetchedPost.Body)

	// Step 5: Upvote the post.
	upvoteBody := `{"author":"user-1","vote":1}`
	req = httptest.NewRequest(http.MethodPut,
		fmt.Sprintf("/servers/%s/posts/%s/vote", createdServer.ID, createdPost.ID),
		strings.NewReader(upvoteBody))
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("upvote post: got status %d, want %d\nbody: %s", w.Code, http.StatusOK, w.Body.String())
	}
	t.Logf("upvoted post %q", createdPost.ID)

	// Step 6: Fetch the post again and verify the vote count.
	req = httptest.NewRequest(http.MethodGet,
		fmt.Sprintf("/servers/%s/posts/%s", createdServer.ID, createdPost.ID), nil)
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("get post after vote: got status %d, want %d\nbody: %s", w.Code, http.StatusOK, w.Body.String())
	}
	var votedPost models.Post
	if err := json.NewDecoder(w.Body).Decode(&votedPost); err != nil {
		t.Fatalf("get post after vote: failed to decode response: %v", err)
	}
	if votedPost.Votes != 1 {
		t.Errorf("get post after vote: got votes %d, want 1", votedPost.Votes)
	}
	t.Logf("confirmed post %q has %d vote(s)", createdPost.ID, votedPost.Votes)
}
