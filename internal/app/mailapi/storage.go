package mailapi

import (
	"sync"

	"github.com/google/uuid"
	"github.com/nilovartem/mail-api/internal/app/model"
)

const (
	NO_LINK = ""
)

// Storage ...
type Storage struct {
	users map[string]*model.User
	mutex sync.Mutex
}

// NewStorage ...
func NewStorage() *Storage {
	return &Storage{
		users: make(map[string]*model.User),
	}
}

// Add ...
func (s *Storage) Add(username string) string {
	s.mutex.Lock()
	link := uuid.New().String()
	u := &model.User{Username: username}
	s.users[link] = u
	s.mutex.Unlock()
	return link
}

// Remove ...
func (s *Storage) Remove(link string) {
	s.mutex.Lock()
	delete(s.users, link)
	s.mutex.Unlock()
}

// GetLink ...
func (s *Storage) GetLink(username string) (string, bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for link, user := range s.users { //TODO: кринж
		if user.Username == username {
			return link, true
		}
	}
	return NO_LINK, false
}

// GetUser ...
func (s *Storage) GetUser(link string) (*model.User, bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if user, found := s.users[link]; found {
		return user, true
	}
	return nil, false
}
