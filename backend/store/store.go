package store

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"
	"github.com/tonitran/dischord/models"
)

type Store struct {
	db *sql.DB
}

// Open opens a Postgres connection, applies the schema, and returns a Store.
func Open(connStr string) (*Store, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}
	if err := ApplySchema(db); err != nil {
		db.Close()
		return nil, err
	}
	return New(db), nil
}

// New wraps an existing *sql.DB. Schema must be applied separately via ApplySchema.
func New(db *sql.DB) *Store {
	return &Store{db: db}
}

// ApplySchema creates all tables if they don't exist.
func ApplySchema(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id         TEXT PRIMARY KEY,
			username   TEXT NOT NULL DEFAULT '',
			email      TEXT NOT NULL DEFAULT '',
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);
		CREATE TABLE IF NOT EXISTS servers (
			id         TEXT PRIMARY KEY,
			name       TEXT NOT NULL DEFAULT '',
			owner_id   TEXT NOT NULL DEFAULT '',
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);
		CREATE TABLE IF NOT EXISTS posts (
			id         TEXT PRIMARY KEY,
			server_id  TEXT NOT NULL DEFAULT '',
			author_id  TEXT NOT NULL DEFAULT '',
			title      TEXT NOT NULL DEFAULT '',
			body       TEXT NOT NULL DEFAULT '',
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);
		CREATE TABLE IF NOT EXISTS votes (
			post_id   TEXT NOT NULL,
			author_id TEXT NOT NULL,
			vote      INTEGER NOT NULL DEFAULT 0,
			PRIMARY KEY (post_id, author_id)
		);
		CREATE TABLE IF NOT EXISTS friends (
			user_id   TEXT NOT NULL,
			friend_id TEXT NOT NULL,
			PRIMARY KEY (user_id, friend_id)
		);
		CREATE TABLE IF NOT EXISTS messages (
			id         TEXT PRIMARY KEY,
			server_id  TEXT NOT NULL DEFAULT '',
			author_id  TEXT NOT NULL DEFAULT '',
			content    TEXT NOT NULL DEFAULT '',
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);
		CREATE TABLE IF NOT EXISTS server_user (
			server_id TEXT NOT NULL REFERENCES servers(id),
			user_id   TEXT NOT NULL REFERENCES users(id),
			PRIMARY KEY (server_id, user_id)
		);
	`)
	return err
}

// TruncateAll removes all rows from every table. Intended for use in tests.
func TruncateAll(db *sql.DB) error {
	_, err := db.Exec(`TRUNCATE TABLE server_user, messages, votes, posts, friends, servers, users`)
	return err
}

// isDuplicateKey reports whether err is a Postgres unique-constraint violation.
func isDuplicateKey(err error) bool {
	var pqErr *pq.Error
	return errors.As(err, &pqErr) && pqErr.Code == "23505"
}

// --- Users ---

func (s *Store) CreateUser(u models.User) error {
	_, err := s.db.Exec(
		`INSERT INTO users (id, username, email, created_at) VALUES ($1, $2, $3, $4)`,
		u.ID, u.Username, u.Email, u.CreatedAt,
	)
	if isDuplicateKey(err) {
		return fmt.Errorf("user %s already exists", u.ID)
	}
	return err
}

func (s *Store) GetUser(id string) (models.User, error) {
	var u models.User
	err := s.db.QueryRow(
		`SELECT id, username, email, created_at FROM users WHERE id = $1`, id,
	).Scan(&u.ID, &u.Username, &u.Email, &u.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return models.User{}, fmt.Errorf("user %s not found", id)
	}
	return u, err
}

// --- Friends ---

func (s *Store) AddFriend(userID, friendID string) error {
	var count int
	if err := s.db.QueryRow(
		`SELECT COUNT(*) FROM users WHERE id = $1 OR id = $2`, userID, friendID,
	).Scan(&count); err != nil {
		return err
	}
	if count < 2 {
		return fmt.Errorf("one or more users not found")
	}
	_, err := s.db.Exec(`
		INSERT INTO friends (user_id, friend_id) VALUES ($1, $2), ($2, $1)
		ON CONFLICT DO NOTHING
	`, userID, friendID)
	return err
}

func (s *Store) GetFriends(userID string) ([]models.User, error) {
	rows, err := s.db.Query(`
		SELECT u.id, u.username, u.email, u.created_at
		FROM users u
		JOIN friends f ON f.friend_id = u.id
		WHERE f.user_id = $1
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var friends []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.CreatedAt); err != nil {
			return nil, err
		}
		friends = append(friends, u)
	}
	return friends, rows.Err()
}

// --- Posts ---

func (s *Store) CreatePost(p models.Post) error {
	_, err := s.db.Exec(
		`INSERT INTO posts (id, server_id, author_id, title, body, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		p.ID, p.ServerID, p.AuthorID, p.Title, p.Body, p.CreatedAt, p.UpdatedAt,
	)
	if isDuplicateKey(err) {
		return fmt.Errorf("post %s already exists", p.ID)
	}
	return err
}

func (s *Store) GetPost(serverID, id string) (models.Post, error) {
	var p models.Post
	err := s.db.QueryRow(`
		SELECT p.id, p.server_id, p.author_id, p.title, p.body,
		       p.created_at, p.updated_at,
		       COALESCE(SUM(v.vote), 0) AS votes
		FROM posts p
		LEFT JOIN votes v ON v.post_id = p.id
		WHERE p.id = $1
		GROUP BY p.id, p.server_id, p.author_id, p.title, p.body, p.created_at, p.updated_at
	`, id).Scan(&p.ID, &p.ServerID, &p.AuthorID, &p.Title, &p.Body, &p.CreatedAt, &p.UpdatedAt, &p.Votes)
	if errors.Is(err, sql.ErrNoRows) {
		return models.Post{}, fmt.Errorf("post %s not found", id)
	}
	return p, err
}

func (s *Store) UpdatePost(p models.Post) error {
	res, err := s.db.Exec(
		`UPDATE posts SET title = $1, body = $2, updated_at = $3 WHERE id = $4`,
		p.Title, p.Body, p.UpdatedAt, p.ID,
	)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return fmt.Errorf("post %s not found", p.ID)
	}
	return nil
}

func (s *Store) DeletePost(id string) error {
	res, err := s.db.Exec(`DELETE FROM posts WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return fmt.Errorf("post %s not found", id)
	}
	return nil
}

func (s *Store) GetVote(postID, authorID string) (models.Vote, error) {
	var v models.Vote
	err := s.db.QueryRow(
		`SELECT post_id, author_id, vote FROM votes WHERE post_id = $1 AND author_id = $2`,
		postID, authorID,
	).Scan(&v.PostID, &v.AuthorID, &v.Vote)
	if errors.Is(err, sql.ErrNoRows) {
		return models.Vote{}, fmt.Errorf("Vote %s not found", postID+"-"+authorID)
	}
	return v, err
}

func (s *Store) PostVote(postID, authorID string, amount int) error {
	var exists bool
	if err := s.db.QueryRow(
		`SELECT EXISTS(SELECT 1 FROM posts WHERE id = $1)`, postID,
	).Scan(&exists); err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("post %s not found", postID)
	}
	_, err := s.db.Exec(`
		INSERT INTO votes (post_id, author_id, vote) VALUES ($1, $2, $3)
		ON CONFLICT (post_id, author_id) DO UPDATE SET vote = EXCLUDED.vote
	`, postID, authorID, amount)
	return err
}

// --- Servers ---

func (s *Store) CreateServer(srv models.Server) error {
	_, err := s.db.Exec(
		`INSERT INTO servers (id, name, owner_id, created_at) VALUES ($1, $2, $3, $4)`,
		srv.ID, srv.Name, srv.OwnerID, srv.CreatedAt,
	)
	if isDuplicateKey(err) {
		return fmt.Errorf("server %s already exists", srv.ID)
	}
	return err
}

func (s *Store) GetServer(id string) (models.Server, error) {
	var srv models.Server
	err := s.db.QueryRow(
		`SELECT id, name, owner_id, created_at FROM servers WHERE id = $1`, id,
	).Scan(&srv.ID, &srv.Name, &srv.OwnerID, &srv.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return models.Server{}, fmt.Errorf("server %s not found", id)
	}
	if err != nil {
		return models.Server{}, err
	}
	rows, err := s.db.Query(`SELECT id FROM posts WHERE server_id = $1`, id)
	if err != nil {
		return models.Server{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var postID string
		if err := rows.Scan(&postID); err != nil {
			return models.Server{}, err
		}
		srv.Posts = append(srv.Posts, postID)
	}
	if err := rows.Err(); err != nil {
		return models.Server{}, err
	}

	memberRows, err := s.db.Query(`SELECT user_id FROM server_user WHERE server_id = $1`, id)
	if err != nil {
		return models.Server{}, err
	}
	defer memberRows.Close()
	for memberRows.Next() {
		var memberID string
		if err := memberRows.Scan(&memberID); err != nil {
			return models.Server{}, err
		}
		srv.MemberIDs = append(srv.MemberIDs, memberID)
	}
	return srv, memberRows.Err()
}

// --- Server Members ---

func (s *Store) JoinServer(serverID, userID string) error {
	var count int
	if err := s.db.QueryRow(
		`SELECT COUNT(*) FROM servers WHERE id = $1`, serverID,
	).Scan(&count); err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("server %s not found", serverID)
	}
	if err := s.db.QueryRow(
		`SELECT COUNT(*) FROM users WHERE id = $1`, userID,
	).Scan(&count); err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("user %s not found", userID)
	}
	_, err := s.db.Exec(
		`INSERT INTO server_user (server_id, user_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`,
		serverID, userID,
	)
	return err
}

func (s *Store) GetServerMembers(serverID string) ([]models.User, error) {
	rows, err := s.db.Query(`
		SELECT u.id, u.username, u.email, u.created_at
		FROM users u
		JOIN server_user su ON su.user_id = u.id
		WHERE su.server_id = $1
	`, serverID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var members []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.CreatedAt); err != nil {
			return nil, err
		}
		members = append(members, u)
	}
	return members, rows.Err()
}

// --- Messages ---

func (s *Store) CreateMessage(m models.Message) error {
	var exists bool
	if err := s.db.QueryRow(
		`SELECT EXISTS(SELECT 1 FROM servers WHERE id = $1)`, m.ServerID,
	).Scan(&exists); err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("server %s not found", m.ServerID)
	}
	_, err := s.db.Exec(
		`INSERT INTO messages (id, server_id, author_id, content, created_at) VALUES ($1, $2, $3, $4, $5)`,
		m.ID, m.ServerID, m.AuthorID, m.Content, m.CreatedAt,
	)
	return err
}

func (s *Store) GetMessagesByServer(serverID string) []models.Message {
	rows, err := s.db.Query(
		`SELECT id, server_id, author_id, content, created_at FROM messages WHERE server_id = $1 ORDER BY created_at`,
		serverID,
	)
	if err != nil {
		return nil
	}
	defer rows.Close()
	var msgs []models.Message
	for rows.Next() {
		var m models.Message
		if err := rows.Scan(&m.ID, &m.ServerID, &m.AuthorID, &m.Content, &m.CreatedAt); err != nil {
			return msgs
		}
		msgs = append(msgs, m)
	}
	return msgs
}
