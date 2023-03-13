package user

import (
	arango "github.com/arangodb/go-driver"
	"github.com/gabriel-ross/trade"
)

type repository struct {
	database arango.Database
}

func (r *repository) Create(user trade.User) (trade.User, error) {
	return trade.User{}, nil
}

func (r *repository) List() ([]trade.User, error) {
	return []trade.User{}, nil
}

func (r *repository) Get(id string) (trade.User, error) {
	return trade.User{}, nil
}

func (r *repository) Update(id string, user trade.User) (trade.User, error) {
	return trade.User{}, nil
}

func (r *repository) Delete(id string) error {
	return nil
}
