package account

import (
	"context"

	arango "github.com/arangodb/go-driver"
	"github.com/gabriel-ross/trade"
)

type repository struct {
	database       arango.Database
	collectionName string
}

func NewRepository(dbClient arango.Database, collectionName string) *repository {
	return &repository{
		database:       dbClient,
		collectionName: collectionName,
	}
}

// Create creates a new document.
func (r *repository) Create(ctx context.Context, account trade.Account) (trade.Account, error) {
	var err error
	col, err := r.database.Collection(ctx, r.collectionName)
	if err != nil {
		return trade.Account{}, err
	}

	meta, err := col.CreateDocument(ctx, account)
	if err != nil {
		return trade.Account{}, err
	}

	account.ID = meta.Key
	return account, nil
}

// List returns all documents in the collection.
func (r *repository) List(ctx context.Context, query string) ([]trade.Account, error) {
	var err error
	results := []trade.Account{}

	cur, err := r.database.Query(ctx, query, map[string]interface{}{})
	if err != nil {
		return []trade.Account{}, err
	}
	defer cur.Close()

	for cur.HasMore() {
		var acc trade.Account
		if _, err = cur.ReadDocument(ctx, &acc); err != nil {
			return []trade.Account{}, err
		}
		results = append(results, acc)
	}

	return results, nil
}

// Get returns account with given id. If no document is found returns NotFoundError.
func (r *repository) Get(ctx context.Context, id string) (trade.Account, error) {
	var err error
	col, err := r.database.Collection(ctx, r.collectionName)
	if err != nil {
		return trade.Account{}, err
	}

	var account trade.Account
	_, err = col.ReadDocument(ctx, id, &account)
	if err != nil {
		return trade.Account{}, err
	}

	return account, nil
}

// Update updates the document identified by id with values account and returns
// the new document. If no document with id is found returns NotFoundError.
func (r *repository) Update(ctx context.Context, id string, account trade.Account) (trade.Account, error) {
	var err error
	col, err := r.database.Collection(ctx, r.collectionName)
	if err != nil {
		return trade.Account{}, err
	}

	var result trade.Account
	_, err = col.UpdateDocument(arango.WithReturnNew(ctx, &result), id, account)
	if err != nil {
		return trade.Account{}, err
	}

	return result, nil
}

// Delete deletes the document with given id from the collection. If no match
// is found returns NotFoundError.
func (r *repository) Delete(ctx context.Context, id string) error {
	var err error
	col, err := r.database.Collection(ctx, r.collectionName)
	if err != nil {
		return err
	}

	_, err = col.RemoveDocument(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
