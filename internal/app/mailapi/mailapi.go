package mailapi

import (
	"fmt"
	"net/http"
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
		s.GetFile(w, r)
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
func (s *Server) GetFile(w http.ResponseWriter, r *http.Request) {
	link := strings.Replace(r.URL.Path, "/", "", -1) //TODO: костыль №1
	link = strings.Replace(link, "get", "", 1)       //TODO: костыль №2
	if u, found := s.storage.FindByLink(link); found {
		fmt.Fprint(w, "\nYour mail is:", u.Mail)
		return
		//zip time
	}
	fmt.Fprint(w, "\nNo zip for you")
	/*
		//fmt.Fprint(w, "Get\n", r.URL.Path)
		link := strings.Replace(r.URL.Path, "/", "", -1)
		link = strings.Replace(link, "get", "", 1)
		//fmt.Fprint(w, "\nLink:", link)
		//if link is exist (and not expired) return ZIP
		if idx := slices.IndexFunc(s.Users, func(user *User) bool { return user.link == link }); idx != -1 {
			mail := s.Users[idx].mail
			w.Header().Set("Content-Type", "application/zip")
			w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s.zip\"", mail))
			zip, err := s.ZipHandler(mail)
			if err != nil {
				s.logger.Errorln(err)
				return
			}
			//w.Write(zip)
			w.Write(zip)
			return
		}
		fmt.Fprint(w, "\nNo zip for you:")
	*/
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
