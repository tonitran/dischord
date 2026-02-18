package store

import (
	"fmt"
	"sync"

	"github.com/tonitran/dischord/models"
)

type Store struct {
	mu       sync.RWMutex
	users    map[string]models.User
	posts    map[string]models.Post
	servers  map[string]models.Server
	messages map[string]models.Message
	friends  map[string]map[string]bool // userID -> set of friendIDs
}

func New() *Store {
	return &Store{
		users:    make(map[string]models.User),
		posts:    make(map[string]models.Post),
		servers:  make(map[string]models.Server),
		messages: make(map[string]models.Message),
		friends:  make(map[string]map[string]bool),
	}
}

// --- Users ---

func (s *Store) CreateUser(u models.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.users[u.ID]; exists {
		return fmt.Errorf("user %s already exists", u.ID)
	}
	s.users[u.ID] = u
	return nil
}

func (s *Store) GetUser(id string) (models.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	u, ok := s.users[id]
	if !ok {
		return models.User{}, fmt.Errorf("user %s not found", id)
	}
	return u, nil
}

// --- Friends ---

func (s *Store) AddFriend(userID, friendID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.users[userID]; !ok {
		return fmt.Errorf("user %s not found", userID)
	}
	if _, ok := s.users[friendID]; !ok {
		return fmt.Errorf("user %s not found", friendID)
	}
	if s.friends[userID] == nil {
		s.friends[userID] = make(map[string]bool)
	}
	if s.friends[friendID] == nil {
		s.friends[friendID] = make(map[string]bool)
	}
	s.friends[userID][friendID] = true
	s.friends[friendID][userID] = true
	return nil
}

func (s *Store) GetFriends(userID string) ([]models.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	friendSet, ok := s.friends[userID]
	if !ok {
		return nil, nil
	}
	var friends []models.User
	for fid := range friendSet {
		if u, ok := s.users[fid]; ok {
			friends = append(friends, u)
		}
	}
	return friends, nil
}

// --- Posts ---

func (s *Store) CreatePost(p models.Post) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.posts[p.ID]; exists {
		return fmt.Errorf("post %s already exists", p.ID)
	}
	s.posts[p.ID] = p
	return nil
}

func (s *Store) GetPost(server_id, id string) (models.Post, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	p, ok := s.posts[id]
	if !ok {
		return models.Post{}, fmt.Errorf("post %s not found", id)
	}
	return p, nil
}

func (s *Store) UpdatePost(p models.Post) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.posts[p.ID]; !ok {
		return fmt.Errorf("post %s not found", p.ID)
	}
	s.posts[p.ID] = p
	return nil
}

func (s *Store) DeletePost(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.posts[id]; !ok {
		return fmt.Errorf("post %s not found", id)
	}
	delete(s.posts, id)
	return nil
}

// --- Servers ---

func (s *Store) CreateServer(srv models.Server) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.servers[srv.ID]; exists {
		return fmt.Errorf("server %s already exists", srv.ID)
	}
	s.servers[srv.ID] = srv
	return nil
}

func (s *Store) GetServer(id string) (models.Server, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	srv, ok := s.servers[id]
	if !ok {
		return models.Server{}, fmt.Errorf("server %s not found", id)
	}
	for postID, post := range s.posts {
		if post.ServerID == id {
			srv.Posts = append(srv.Posts, postID)
		}
	}
	return srv, nil
}

// --- Messages ---

func (s *Store) CreateMessage(m models.Message) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.servers[m.ServerID]; !ok {
		return fmt.Errorf("server %s not found", m.ServerID)
	}
	s.messages[m.ID] = m
	return nil
}

func (s *Store) GetMessagesByServer(serverID string) []models.Message {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var msgs []models.Message
	for _, m := range s.messages {
		if m.ServerID == serverID {
			msgs = append(msgs, m)
		}
	}
	return msgs
}
