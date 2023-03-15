package app

import (
	"log"
	"net/http"

	arango "github.com/arangodb/go-driver"
	arangohttp "github.com/arangodb/go-driver/http"
	"github.com/gabriel-ross/trade/user"
	"github.com/go-chi/chi"
)

// Config contains all the settings for an application instance.
type Config struct {
	PORT             string `env:"PORT" default:"8080" required:"false"`
	DB_ADDRESS       string `env:"DB_ADDRESS" default:"8080" required:"true"`
	DB_NAME          string `env:"DB_NAME" default:"demo" required:"false"`
	createOnNotExist bool
}

// application is the entrypoint to the program and houses the necessary
// dependencies.
type application struct {
	cnf      Config
	router   chi.Router
	dbClient arango.Database
}

func New(cnf Config, options ...func(*application)) *application {
	a := &application{
		cnf:    cnf,
		router: chi.NewRouter(),
	}

	// Configure options
	for _, option := range options {
		option(a)
	}

	// Instantiate ArangoDB connection
	arangoConn, err := arangohttp.NewConnection(arangohttp.ConnectionConfig{
		Endpoints: []string{a.cnf.DB_ADDRESS},
	})
	if err != nil {
		log.Fatalf("error connecting to db: %v", err)
	}

	arangoClient, err := arango.NewClient(arango.ClientConfig{
		Connection: arangoConn,
	})
	if err != nil {
		log.Fatalf("error instantiating arango client: %v", err)
	}

	dbClient, err := arangoClient.Database(nil, a.cnf.DB_NAME)
	if err != nil {
		if arango.IsNotFoundGeneral(err) && a.cnf.createOnNotExist {
			dbClient, err = arangoClient.CreateDatabase(nil, a.cnf.DB_NAME, nil)
			if err != nil {
				log.Fatalf("error creating database: %v", err)
			}
		} else {
			log.Fatalf("error connecting to database: %v", err)
		}
	}
	a.dbClient = dbClient

	// Instantiate and register services
	user.New(a.router, "/users", user.NewUserRepository(a.dbClient, "users"))

	return a
}

func WithCreateOnNotExist(flag bool) func(*application) {
	return func(a *application) {
		a.cnf.createOnNotExist = flag
	}
}

func (a *application) Run() error {
	return http.ListenAndServe(":"+a.cnf.PORT, a.router)
}
