package models

import "time"

type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type FriendRequest struct {
	UserID   string `json:"user_id"`
	FriendID string `json:"friend_id"`
}

type Post struct {
	ID        string    `json:"id"`
	ServerID  string    `json:"server_id"`
	AuthorID  string    `json:"author_id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	Votes     int       `json:"votes"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Vote struct {
	PostID   string `json:"post_id"`
	AuthorID string `json:"author_id"`
	Vote     int    `json:"vote"`
}

type Server struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	OwnerID   string    `json:"owner_id"`
	MemberIDs []string  `json:"member_ids"`
	Posts     []string  `json:"post_ids"`
	CreatedAt time.Time `json:"created_at"`
}

type Message struct {
	ID        string    `json:"id"`
	ServerID  string    `json:"server_id"`
	AuthorID  string    `json:"author_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}
