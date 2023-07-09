package mailapi

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/nilovartem/mail-api/internal/app/config"
	"github.com/sirupsen/logrus"
)

// Server ...
type Server struct {
	config  *config.Config
	logger  *logrus.Logger
	storage *Storage
}

// NewServer ...
func NewServer(c *config.Config) *Server {
	return &Server{
		config:  c,
		logger:  logrus.New(),
		storage: NewStorage(),
	}
}

// Start ...
func (s *Server) Start() error {
	s.logger.Infof("Starting server on %s\n", s.config.BindAddress)
	if err := s.configureLogger(); err != nil {
		return err
	}
	mux := http.NewServeMux()

	PostHandler := http.HandlerFunc(s.postHandler)
	mux.Handle("/", s.serveHTTP(s.staticAuth(PostHandler)))

	GetHandler := http.HandlerFunc(s.getHandler)
	mux.Handle("/get/", s.serveHTTP(GetHandler))

	return http.ListenAndServe(s.config.BindAddress, mux)
}

var (
	postHandlerRe = regexp.MustCompile(`^\/[^\/]+$`)
	getHandlerRe  = regexp.MustCompile(`^\/get\/[^\/]+$`)
)

// serveHTTP ...
func (s *Server) serveHTTP(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && postHandlerRe.MatchString(r.URL.Path):
			next.ServeHTTP(w, r)
		case r.Method == http.MethodGet && getHandlerRe.MatchString(r.URL.Path):
			next.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// staticAuth ...
func (s *Server) staticAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		usernameURL := strings.Trim(r.URL.Path, "/")
		username, password, ok := r.BasicAuth()
		if ok && username == usernameURL {
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
	username := strings.Trim(r.URL.Path, "/")
	link, found := s.storage.GetLink(username)
	if !found {
		link = s.storage.Add(username)
		go func() {
			ticker := time.NewTicker(s.config.TTL)
			for {
				<-ticker.C
				s.storage.Remove(link)
			}
		}()
	}
	fmt.Fprint(w, "/get/"+link)
}

// getHandler returns zip by UUID
func (s *Server) getHandler(w http.ResponseWriter, r *http.Request) {
	link := strings.Replace(r.URL.Path, "/get/", "", -1)
	if u, found := s.storage.GetUser(link); found {
		zip, err := u.Zip(filepath.Join(s.config.Mailbox, u.Username), u.Username, s.config.PDF)
		if err != nil {
			http.Error(w, "zip error", http.StatusInternalServerError)
			s.logger.Errorln(err)
			return
		}
		w.Header().Set("Content-Type", "application/zip")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.zip", u.Username))
		w.Write(zip)
		return
	}
	http.Error(w, "link expired", http.StatusBadRequest)
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
