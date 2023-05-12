package user

import "github.com/go-chi/chi"

// Routes returns a new chi router with all user routes mounted to it.
func (s *service) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/", s.handleCreate())
	r.Get("/", s.handleList())
	r.Route("/{id}", func(r chi.Router) {
		r.Get("/", s.handleGet())
		r.Put("/", s.handlePut())
		r.Delete("/", s.handleDelete())
		r.Get("/accounts", s.handleGetAccounts())
	})

	return r
}
