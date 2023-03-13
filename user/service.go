package user

import (
	"github.com/gabriel-ross/trade"
	"github.com/go-chi/chi"
)

// Repository is the API for the User datastore.
type Repository interface {
	Create(trade.User) (trade.User, error)
	List() ([]trade.User, error)
	Get(string) (trade.User, error)
	Update(string, trade.User) (trade.User, error)
	Delete(string) error
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
