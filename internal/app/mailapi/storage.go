package mailapi

import (
	"fmt"
	"os"
	"sync"

	"github.com/nilovartem/mail-api/internal/app/model"
	"golang.org/x/exp/slices"
)

// Storage ...
type Storage struct {
	users  []*model.User // User{username, link}
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
func (s *Storage) Inflate() error {
	entries, err := os.ReadDir(s.config.Mailbox)
	if err != nil {
		return err
	}
	for _, e := range entries {
		s.users = append(s.users, &model.User{Mail: e.Name(), Link: ""})
	}
	return nil
}

// FindByMail ...
func (s *Storage) FindByMail(mail string) (*model.User, bool, error) {
	//var link string
	//if YOUR user is not exists
	//not found
	var id int
	if id = slices.IndexFunc(s.users, func(user *model.User) bool { return user.Mail == mail }); id == -1 {
		return nil, false, fmt.Errorf("user not found")
	}
	//check if user exists and link is also exists
	if idx := slices.IndexFunc(s.users, func(user *model.User) bool { return user.Mail == mail && user.Link != "" }); idx != -1 {
		return s.users[idx], true, nil
	}
	//default - user exists, link not
	return s.users[id], false, nil
}

// FindByLink
func (s *Storage) FindByLink(link string) (*model.User, bool) {
	var id int
	if id = slices.IndexFunc(s.users, func(user *model.User) bool { return user.Link == link }); id != -1 {
		return s.users[id], true
	}
	return nil, false
}
