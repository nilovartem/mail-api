package mailapi

import (
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
	if err := s.configureLogger(); err != nil {
		return err
	}
	//inflate storage with users (emails)
	err := s.storage.Inflate()
	if err != nil {
		return err
	}
	s.logger.Infoln("starting server")
	mux := http.NewServeMux()
	mux.Handle("/", s)
	mux.Handle("/get/", s)
	return http.ListenAndServe(s.config.BindAddress, mux)
}

var (
	getLinkRe = regexp.MustCompile(`^\/[^\/]+$`)
	getFileRe = regexp.MustCompile(`^\/get\/[^\/]+$`)
)

// ServeHTTP ...
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodPost && getLinkRe.MatchString(r.URL.Path):
		err := s.GetLink(w, r)
		if err != nil {
			s.logger.Errorln(err)
		}
		return
	case r.Method == http.MethodGet && getFileRe.MatchString(r.URL.Path):
		err := s.GetFile(w, r)
		if err != nil {
			s.logger.Errorln(err)
		}
		return
	default:
		http.NotFound(w, r) //TODO: maybe if not allowed or if not match then not found
		return
	}
}

// GetLink return UUID....
func (s *Server) GetLink(w http.ResponseWriter, r *http.Request) error {
	fmt.Fprint(w, "GetLink\n", r.URL.Path)
	mail := strings.Trim(r.URL.Path, "/")
	fmt.Fprint(w, "\nmail:", mail)
	u, hasLink, err := s.storage.FindByMail(mail)
	if err != nil {
		fmt.Fprint(w, "\nUser not found")
		return err
	}
	if hasLink {
		link := u.Link
		fmt.Fprint(w, "\nOld link:", link)
		return nil
	}
	//if not found
	u.NewLink(s.config.TTL)
	fmt.Fprint(w, "\nNew link:", u.Link)
	s.logger.Info(s.storage.users)
	return nil
}

// GetZip returns zip by UUID
func (s *Server) GetFile(w http.ResponseWriter, r *http.Request) error {
	link := strings.Replace(r.URL.Path, "/", "", -1) //TODO: костыль №1
	link = strings.Replace(link, "get", "", 1)       //TODO: костыль №2
	if u, found := s.storage.FindByLink(link); found {
		s.logger.Infoln("\nYour mail is:", u.Mail)
		//TODO: zip time
		zip, err := u.Zip(filepath.Join(s.config.Mailbox, u.Mail), u.Mail)
		if err != nil {
			s.logger.Errorln(err)
			return err
		}
		s.logger.Infoln("buffer size if ", len(zip))
		w.Header().Set("Content-Type", "application/zip")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.zip", u.Mail))
		w.Write(zip)
		return nil
	}
	fmt.Fprint(w, "\nNo zip for you")
	return nil
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
