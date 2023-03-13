package trade

import (
	arango "github.com/arangodb/go-driver"
	"github.com/go-chi/chi"
)

// Application is the entrypoint to the application and houses the necessary
// dependencies.
type Application struct {
	router   chi.Router
	database arango.Database
}
