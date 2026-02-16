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
	servers := &handlers.ServerHandler{Store: s}
	messages := &handlers.MessageHandler{Store: s}

	// Users
	mux.HandleFunc("POST /users", users.Create)
	mux.HandleFunc("GET /users/{id}", users.Get)

	// Friends
	mux.HandleFunc("POST /users/{id}/friends", friends.Add)
	mux.HandleFunc("GET /users/{id}/friends", friends.List)

	// Posts
	mux.HandleFunc("POST /posts", posts.Create)
	mux.HandleFunc("GET /posts/{id}", posts.Get)
	mux.HandleFunc("PATCH /posts/{id}", posts.Update)
	mux.HandleFunc("DELETE /posts/{id}", posts.Delete)

	// Servers
	mux.HandleFunc("POST /servers", servers.Create)
	mux.HandleFunc("GET /servers/{id}", servers.Get)

	// Messages
	mux.HandleFunc("POST /servers/{server_id}/messages", messages.Create)
	mux.HandleFunc("GET /servers/{server_id}/messages", messages.ListByServer)

	return mux
}
