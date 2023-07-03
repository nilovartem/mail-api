package mailapi

// Mailapi ...
type Server struct {
	config *Config
}

// New ...
func NewServer(config *Config) *Server {
	return &Server{
		config: config,
	}
}

// Start ...
func (server *Server) Start() error {
	return nil
}
