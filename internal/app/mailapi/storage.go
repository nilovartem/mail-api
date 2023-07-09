package mailapi

import (
	"sync"

	"github.com/nilovartem/mail-api/internal/app/model"
	"golang.org/x/exp/slices"
)

// Storage ...
type Storage struct {
	users  []*model.User
	mutex  sync.Mutex
	config *Config
}

// NewStorage ...
func NewStorage(c *Config) *Storage {
	return &Storage{
		users:  []*model.User{},
		mutex:  sync.Mutex{},
		config: c,
	}
}

// Inflate ...
func (s *Storage) Inflate(users *map[string][32]byte) {
	for u := range *users {
		s.users = append(s.users, &model.User{Mail: u, Link: model.NO_LINK})
	}
}

// FindByMail ...
func (s *Storage) FindByMail(mail string) (*model.User, bool) {
	id := slices.IndexFunc(s.users, func(user *model.User) bool { return user.Mail == mail })
	if s.users[id].Link != model.NO_LINK {
		return s.users[id], true
	}
	return s.users[id], false
}

// FindByLink
func (s *Storage) FindByLink(link string) (*model.User, bool) {
	var id int
	if id = slices.IndexFunc(s.users, func(user *model.User) bool { return user.Link == link }); id != -1 {
		return s.users[id], true
	}
	return nil, false
}
