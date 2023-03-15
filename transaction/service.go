package transaction

import (
	"context"

	"github.com/gabriel-ross/trade"
	"github.com/go-chi/chi"
)

// Repository is the API for the Transaction datastore.
type Repository interface {
	Create(ctx context.Context, t trade.Transaction) (trade.Transaction, error)
	List(ctx context.Context) ([]trade.Transaction, error)
	Get(ctx context.Context, id string) (trade.Transaction, error)
	Update(ctx context.Context, id string, t trade.Transaction) (trade.Transaction, error)
	Delete(ctx context.Context, id string) error
}

// Service houses the API and necessary dependencies for interacting with
// account resources.
type service struct {
	router   chi.Router
	database Repository
}

// New mounts the account routes on r at endpoint and returns a new account service.
func New(r chi.Router, endpoint string, database Repository, options ...func(*service)) *service {
	svc := &service{
		router:   r,
		database: database,
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
