package user

import (
	"context"

	"github.com/gabriel-ross/trade"
	"github.com/go-chi/chi"
)

// Repository is the API for the User datastore.
type Repository interface {
	Create(ctx context.Context, u trade.User) (trade.User, error)
	List(ctx context.Context) ([]trade.User, error)
	Get(ctx context.Context, id string) (trade.User, error)
	Update(ctx context.Context, id string, u trade.User) (trade.User, error)
	Delete(ctx context.Context, id string) error
}

// Service houses the API and necessary dependencies for interacting with user
// resources.
type service struct {
	router   chi.Router
	database Repository
}

// New mounts the user routes on r at endpoint and returns a new user service.
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

// WithRepository is a functional option for configuring a user service's
// repository upon instantiation.
func WithRepository(repo Repository) func(*service) {
	return func(s *service) {
		s.database = repo
	}
}
