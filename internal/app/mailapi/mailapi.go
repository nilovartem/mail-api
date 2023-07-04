package mailapi

import (
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// Link ...
type Link struct {
	token string
	ttl   time.Duration
}

// Mailapi ...
type Server struct {
	config *Config
	logger *logrus.Logger
	users  map[string]Link // username.lastname@wb.work = {uuid, 1s}
	mutex  sync.Mutex
}

// New ...
func NewServer(c *Config) *Server {
	return &Server{
		config: c,
		logger: logrus.New(),
	}
}

// Start ...
func (s *Server) Start() error {
	if err := s.configureLogger(); err != nil {
		return err
	}
	s.logger.Infoln("starting server")
	return nil
}
func (s *Server) configureLogger() error {
	level, err := logrus.ParseLevel(s.config.LogLevel)
	if err != nil {
		return err
	}
	s.logger.Level = level
	return nil
}
