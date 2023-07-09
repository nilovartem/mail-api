package mailapi

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"
)

// Server ...
type Server struct {
	config  *Config
	logger  *logrus.Logger
	storage *Storage
}

// New ...
func NewServer(c *Config) *Server {
	return &Server{
		config:  c,
		logger:  logrus.New(),
		storage: NewStorage(c),
	}
}

// Start ...
func (s *Server) Start() error {
	s.logger.Infof("Starting server on %s\n", s.config.BindAddress)
	if err := s.configureLogger(); err != nil {
		return err
	}
	s.storage.Inflate(&s.config.Users)

	mux := http.NewServeMux()

	PostHandler := http.HandlerFunc(s.postHandler)
	mux.Handle("/", s.serveHTTP(s.auth(PostHandler)))

	GetHandler := http.HandlerFunc(s.getHandler)
	mux.Handle("/get/", s.serveHTTP(GetHandler))

	return http.ListenAndServe(s.config.BindAddress, mux)
}

var (
	getLinkRe = regexp.MustCompile(`^\/[^\/]+$`)
	getFileRe = regexp.MustCompile(`^\/get\/[^\/]+$`)
)

// ServeHTTP ...
func (s *Server) serveHTTP(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && getLinkRe.MatchString(r.URL.Path):
			next.ServeHTTP(w, r)
		case r.Method == http.MethodGet && getFileRe.MatchString(r.URL.Path):
			next.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// auth ...
func (s *Server) auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		usernamePath := strings.Trim(r.URL.Path, "/")
		username, password, ok := r.BasicAuth()
		if ok && username == usernamePath {
			if passwordHash, found := s.config.Users[username]; found {
				if passwordHash == sha256.Sum256([]byte(password)) {
					next.ServeHTTP(w, r)
					return
				}
			}
		}
		w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	})
}

// postHandler returns UUID....
func (s *Server) postHandler(w http.ResponseWriter, r *http.Request) {
	mail := strings.Trim(r.URL.Path, "/")
	u, hasLink := s.storage.FindByMail(mail)
	if !hasLink {
		u.NewLink(s.config.TTL)
	}
	fmt.Fprint(w, s.config.BindAddress+"/get/"+u.Link)
}

// getHandler returns zip by UUID
func (s *Server) getHandler(w http.ResponseWriter, r *http.Request) {
	link := strings.Replace(r.URL.Path, "/get/", "", -1)
	if u, found := s.storage.FindByLink(link); found {
		zip, err := u.Zip(filepath.Join(s.config.Mailbox, u.Mail), u.Mail)
		if err != nil {
			s.logger.Errorln(err)
			return
		}
		w.Header().Set("Content-Type", "application/zip")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.zip", u.Mail))
		w.Write(zip)
		return
	}
	fmt.Fprint(w, "No zip for you")
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
