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

type Server struct {
	config  *config.Config
	logger  *logrus.Logger
	storage *Storage
}

func NewServer(c *config.Config) *Server {
	return &Server{
		config:  c,
		logger:  logrus.New(),
		storage: NewStorage(),
	}
}

// Start function starts listening and implements the router
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

// regexp for handle routes
var (
	postHandlerRe = regexp.MustCompile(`^\/[^\/]+$`)
	getHandlerRe  = regexp.MustCompile(`^\/get\/[^\/]+$`)
)

// serveHTTP dispathes requests between handlers...
func (s *Server) serveHTTP(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && postHandlerRe.MatchString(r.URL.Path):
			s.logger.Infoln("[serveHTTP] got request to endpoint #1")
			next.ServeHTTP(w, r)
		case r.Method == http.MethodGet && getHandlerRe.MatchString(r.URL.Path):
			s.logger.Infoln("[serveHTTP] got request to endpoint #2")
			next.ServeHTTP(w, r)
		default:
			s.logger.Infoln("[serveHTTP] got invalid request")
			http.NotFound(w, r)
		}
	})
}

// staticAuth checks credentials and passes hanldle functions to postHandler
func (s *Server) staticAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		usernameURL := strings.Trim(r.URL.Path, "/")
		username, password, ok := r.BasicAuth()
		s.logger.Infof("[staticAuth] username: %s, password : %s, username in URL %s\n", username, password, usernameURL)
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

// postHandler responds with link - old or new
func (s *Server) postHandler(w http.ResponseWriter, r *http.Request) {
	username := strings.Trim(r.URL.Path, "/")
	s.logger.Infof("[postHandler] username: %s\n", username)
	link, found := s.storage.GetLink(username)
	if !found {
		s.logger.Infoln("[postHandler] link not found, creating new link")
		link = s.storage.Add(username)
		go func() {
			ticker := time.NewTicker(s.config.TTL)
			for {
				<-ticker.C
				s.storage.Remove(link)
			}
		}()
	}
	s.logger.Infof("[postHandler] link: %s\n", link)
	fmt.Fprint(w, "/get/"+link)
}

// getHandler responds to valid link with zip
func (s *Server) getHandler(w http.ResponseWriter, r *http.Request) {
	link := strings.Replace(r.URL.Path, "/get/", "", -1)
	s.logger.Infof("[getHandler] link: %s\n", link)
	if u, found := s.storage.GetUser(link); found {
		s.logger.Infof("[getHandler] found owner's username: %s\n", u.Username)
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
	s.logger.Infof("[getHandler] link %s is expired\n", link)
	http.Error(w, "link expired", http.StatusBadRequest)
}

// configureLogger sets logLevel
func (s *Server) configureLogger() error {
	level, err := logrus.ParseLevel(s.config.LogLevel)
	if err != nil {
		return err
	}
	s.logger.Level = level
	return nil
}
