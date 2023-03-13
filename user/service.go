package user

import (
	"context"

	"github.com/gabriel-ross/trade"
	"github.com/go-chi/chi"
)

// Repository is the API for the User datastore.
type Repository interface {
	Create(context.Context, trade.User) (trade.User, error)
	List(context.Context) ([]trade.User, error)
	Get(context.Context, string) (trade.User, error)
	Update(context.Context, string, trade.User) (trade.User, error)
	Delete(context.Context, string) error
}

// Service houses the API and necessary dependencies for interacting with user
// resources.
type service struct {
	router   chi.Router
	database Repository
}

// New returns a new user service.
func New() *service {
	return &service{}
}
