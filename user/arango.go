package user

import (
	"context"
	"fmt"

	arango "github.com/arangodb/go-driver"
	"github.com/gabriel-ross/trade"
)

type userRepository struct {
	database           arango.Database
	userCollectionName string
}

func NewUserRepository(dbClient arango.Database, userCollectionName string) *userRepository {
	return &userRepository{
		database:           dbClient,
		userCollectionName: userCollectionName,
	}
}

// Create creates a new document.
func (r *userRepository) Create(ctx context.Context, user trade.User) (trade.User, error) {
	var err error
	col, err := r.database.Collection(ctx, r.userCollectionName)
	if err != nil {
		return trade.User{}, err
	}

	meta, err := col.CreateDocument(ctx, user)
	if err != nil {
		return trade.User{}, err
	}

	user.ID = meta.Key
	return user, nil
}

// List returns all documents in the collection.
func (r *userRepository) List(ctx context.Context) ([]trade.User, error) {
	var err error
	var results []trade.User

	query := fmt.Sprintf(
		`FOR entry IN %s
RETURN entry`, r.userCollectionName)
	cur, err := r.database.Query(ctx, query, map[string]interface{}{})
	if err != nil {
		return []trade.User{}, err
	}
	defer cur.Close()

	for cur.HasMore() {
		var user trade.User
		if _, err = cur.ReadDocument(ctx, &user); err != nil {
			return []trade.User{}, err
		}
		results = append(results, user)
	}

	return results, nil
}

// Get returns user with given id. If no document is found returns NotFoundError.
func (r *userRepository) Get(ctx context.Context, id string) (trade.User, error) {
	var err error
	col, err := r.database.Collection(ctx, r.userCollectionName)
	if err != nil {
		return trade.User{}, err
	}

	var user trade.User
	_, err = col.ReadDocument(ctx, id, &user)
	if err != nil {
		return trade.User{}, err
	}

	return user, nil
}

// Update updates the document identified by id with values user and returns
// the new document. If no document with id is found returns NotFoundError.
func (r *userRepository) Update(ctx context.Context, id string, user trade.User) (trade.User, error) {
	var err error
	col, err := r.database.Collection(ctx, r.userCollectionName)
	if err != nil {
		return trade.User{}, err
	}

	var result trade.User
	_, err = col.UpdateDocument(arango.WithReturnNew(ctx, &result), id, user)
	if err != nil {
		return trade.User{}, err
	}

	return result, nil
}

// Delete deletes the document with given id from the collection. If no match
// is found returns NotFoundError.
func (r *userRepository) Delete(ctx context.Context, id string) error {
	var err error
	col, err := r.database.Collection(ctx, r.userCollectionName)
	if err != nil {
		return err
	}

	_, err = col.RemoveDocument(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
