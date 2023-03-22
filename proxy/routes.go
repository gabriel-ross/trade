package proxy

import (
	"github.com/go-chi/chi"
)

func (s *Server) Routes() chi.Router {
	r := chi.NewRouter()

	return r
}
