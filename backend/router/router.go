package router

import (
	"net/http"

	"github.com/tonitran/dischord/handlers"
	"github.com/tonitran/dischord/store"
)

func New(s *store.Store) http.Handler {
	mux := http.NewServeMux()

	users := &handlers.UserHandler{Store: s}
	friends := &handlers.FriendHandler{Store: s}
	posts := &handlers.PostHandler{Store: s}
	votes := &handlers.VoteHandler{Store: s}
	servers := &handlers.ServerHandler{Store: s}
	messages := &handlers.MessageHandler{Store: s}

	// Servers
	mux.HandleFunc("POST /servers", servers.Create)
	mux.HandleFunc("GET /servers/{id}", servers.Get)

	// Posts
	mux.HandleFunc("POST /servers/{server_id}/posts", posts.Create)
	mux.HandleFunc("GET /servers/{server_id}/posts/{id}", posts.Get)
	mux.HandleFunc("PUT /servers/{server_id}/posts/{id}", posts.Update)
	mux.HandleFunc("DELETE /servers/{server_id}/posts/{id}", posts.Delete)

	// Votes
	mux.HandleFunc("GET /servers/{server_id}/posts/{id}/vote", votes.GetVote)
	mux.HandleFunc("PUT /servers/{server_id}/posts/{id}/vote", votes.PutVote)

	// Users
	mux.HandleFunc("POST /users", users.Create)
	mux.HandleFunc("GET /users/{id}", users.Get)

	// Friends
	mux.HandleFunc("POST /users/{id}/friends", friends.Add)
	mux.HandleFunc("GET /users/{id}/friends", friends.List)

	// Messages
	mux.HandleFunc("POST /servers/{server_id}/messages", messages.Create)
	mux.HandleFunc("GET /servers/{server_id}/messages", messages.ListByServer)

	return mux
}
