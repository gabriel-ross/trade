package account

import (
	"context"
	"net/http"

	"github.com/gabriel-ross/trade"
	"github.com/go-chi/chi"
)

// Repository is the API for the Account datastore.
type Repository interface {
	Create(ctx context.Context, a trade.Account) (string, trade.Account, error)
	Query(ctx context.Context, query string) ([]trade.Account, error)
	Get(ctx context.Context, id string) (trade.Account, error)
	Update(ctx context.Context, id string, a trade.Account) (trade.Account, error)
	Delete(ctx context.Context, id string) error
}

type Renderer interface {
	RenderJSON(w http.ResponseWriter, r *http.Request, httpStatusCode int, body interface{})
	RenderError(w http.ResponseWriter, r *http.Request, svrErr error, code int, format string, args ...any)
}

// Service houses the API and necessary dependencies for interacting with
// account resources.
type service struct {
	router   chi.Router
	database Repository
	renderer Renderer
}

// New mounts the account routes on r at endpoint and returns a new account service.
func New(r chi.Router, endpoint string, database Repository, renderer Renderer, options ...func(*service)) *service {
	svc := &service{
		router:   r,
		database: database,
		renderer: renderer,
	}
	r.Mount(endpoint, svc.Routes())

	for _, option := range options {
		option(svc)
	}

	return svc
}

// WithRepository is a functional option for configuring a account service's
// repository upon instantiation.
func WithRepository(repo Repository) func(*service) {
	return func(s *service) {
		s.database = repo
	}
}
