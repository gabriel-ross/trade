package account

import (
	"context"
	"net/http"

	arango "github.com/arangodb/go-driver"
	"github.com/gabriel-ross/trade"
	"github.com/go-chi/chi"
)

// Repository is the API for the Account datastore.
type Repository interface {
	Create(ctx context.Context, a trade.Account) (trade.Account, error)
	List(ctx context.Context) ([]trade.Account, error)
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
func New(r chi.Router, endpoint string, db arango.Database, collectionName string, renderer Renderer) *service {
	svc := &service{
		router:   r,
		database: NewRepository(db, collectionName),
		renderer: renderer,
	}
	r.Mount(endpoint, svc.Routes())

	return svc
}
