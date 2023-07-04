package mailapi

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"
)

type User struct {
	username string
	link     string
}

// Mailapi ...
type Server struct {
	config *Config
	logger *logrus.Logger
	Users  []User // User{username, link}
	mutex  sync.Mutex
}

// New ...
func NewServer(c *Config) *Server {
	return &Server{
		config: c,
		logger: logrus.New(),
		Users:  []User{},
		mutex:  sync.Mutex{},
	}
}

var (
	authRe   = regexp.MustCompile(`^\/[^\/]+$`)
	getZipRe = regexp.MustCompile(`^\/get\/[^\/]+$`)
)

// ServeHTTP ...
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("content-type", "application/json")
	switch {
	case r.Method == http.MethodPost && authRe.MatchString(r.URL.Path):
		s.Auth(w, r)
		return
	case r.Method == http.MethodGet && getZipRe.MatchString(r.URL.Path):
		s.GetZip(w, r)
		return
	default:
		http.NotFound(w, r) //TODO: maybe if not allowed or if not match then not found
		return
	}

}

// Auth return UUID....
func (s *Server) Auth(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Auth\n", r.URL.Path)
	username := strings.Trim(r.URL.Path, "/")
	fmt.Fprint(w, "\nUsername:", username)

	var link string
	if idx := slices.IndexFunc(s.Users, func(user User) bool { return user.username == username }); idx != -1 {
		link = s.Users[idx].link
		fmt.Fprint(w, "\nYour link:", link)
		return
	}
	link = uuid.New().String()
	s.Users = append(s.Users, User{username: username, link: link})
	fmt.Fprint(w, "\nYour link:", link)
	//start timer
	s.logger.Info(s.Users)
}

// GetZip returns zip by UUID
func (s *Server) GetZip(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Get\n", r.URL.Path)
	link := strings.Replace(r.URL.Path, "/", "", -1)
	link = strings.Replace(link, "get", "", 1)
	fmt.Fprint(w, "\nLink:", link)
	//if link is exist (and not expired) return ZIP
	if idx := slices.IndexFunc(s.Users, func(user User) bool { return user.link == link }); idx != -1 {
		username := s.Users[idx].username
		fmt.Fprint(w, "\nYour ZIP:", username)
		return
	}
	fmt.Fprint(w, "\nNo zip for you:")
}

// Start ...
func (s *Server) Start() error {
	if err := s.configureLogger(); err != nil {
		return err
	}
	s.logger.Infoln("starting server")
	mux := http.NewServeMux()
	mux.Handle("/", s)
	mux.Handle("/get/", s)
	return http.ListenAndServe(s.config.BindAddress, mux)
}

// configureLogger ...
func (s *Server) configureLogger() error {
	level, err := logrus.ParseLevel(s.config.LogLevel)
	if err != nil {
		return err
	}
	s.logger.Level = level
	return nil
}
