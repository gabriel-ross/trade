package app

import (
	"context"
	"fmt"
	"log"
	"net/http"

	arangodriver "github.com/arangodb/go-driver"
	"github.com/gabriel-ross/trade"
	"github.com/gabriel-ross/trade/account"
	"github.com/gabriel-ross/trade/transaction"
	"github.com/gabriel-ross/trade/user"
	"github.com/go-chi/chi"
)

// Config contains all the settings for an application instance.
type Config struct {
	PORT             string `env:"PORT" default:"80" required:"false"`
	DB_ADDRESS       string `env:"DB_ADDRESS" default:"8529" required:"true"`
	DB_NAME          string `env:"DB_NAME" default:"demo" required:"false"`
	createOnNotExist bool
}

// application is the entrypoint to the program and houses the necessary
// dependencies.
type application struct {
	cnf      Config
	router   chi.Router
	dbClient arangodriver.Database
}

// New instantiates a new application according to cnf and options and returns
// the new application.
func New(cnf Config, options ...func(*application)) *application {
	a := &application{
		cnf:    cnf,
		router: chi.NewRouter(),
	}

	// Configure options
	for _, option := range options {
		option(a)
	}

	arangoClient, err := trade.NewArangoClient([]string{a.cnf.DB_ADDRESS})
	if err != nil {
		log.Fatalf("error instantiating arangodb client %v", err)
	}

	a.dbClient, err = arangoClient.Database(context.TODO(), a.cnf.DB_NAME, true, "./db/arango_schema.json")
	if err != nil {
		log.Fatalf("error connecting to database %v", err)
	}

	a.router.Get("/ping", a.Ping())

	// Instantiate and register services
	user.New(a.router, "/users", trade.NewArangoRepository[trade.User](a.dbClient, "users"), &trade.RenderService{})
	account.New(a.router, "/accounts", trade.NewArangoRepository[trade.Account](a.dbClient, "accounts"), &trade.RenderService{})
	transaction.New(a.router, "/transactions", trade.NewArangoRepository[trade.Transaction](a.dbClient, "transactions"), &trade.RenderService{})

	return a
}

// WithCreateOnNotExist is an application functional option. If set to true
// when the application is instantiated if no database with a.cnf.DB_NAME is
// found a database with this name will be created along with any required
// database connections.
func WithCreateOnNotExist(flag bool) func(*application) {
	return func(a *application) {
		a.cnf.createOnNotExist = flag
	}
}

// Run runs the application on a.cnf.PORT
func (a *application) Run() error {
	fmt.Println("application running on port ", a.cnf.PORT)
	return http.ListenAndServe(":"+a.cnf.PORT, a.router)
}

func (a *application) Ping() http.HandlerFunc {
	resp := fmt.Sprintf("server is healthy and running at port %s", a.cnf.PORT)
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(resp))
	}
}
