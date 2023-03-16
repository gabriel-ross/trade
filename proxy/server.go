package proxy

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
)

type Config struct {
	NAME           string
	PORT           string `env:"PORT" default:"8081" required:"false"`
	SERVER_ADDRESS string `env:"SERVER_ADDRESS" default:"localhost:8080" required:"false"`
	CACHE_TIMEOUT  time.Duration
}

type Server struct {
	cnf    Config
	router chi.Router
}

func NewServer(cnf Config) *Server {
	return &Server{
		cnf:    cnf,
		router: chi.NewRouter(),
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Server) Run() error {
	log.Printf("%s server running on port %s", s.cnf.NAME, s.cnf.PORT)
	return http.ListenAndServe(":"+s.cnf.PORT, s)
}
