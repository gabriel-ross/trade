package transaction

import (
	"context"
	"fmt"

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
func (r *repository) Create(ctx context.Context, transaction trade.Transaction) (trade.Transaction, error) {
	var err error
	col, err := r.database.Collection(ctx, r.collectionName)
	if err != nil {
		return trade.Transaction{}, err
	}

	meta, err := col.CreateDocument(ctx, transaction)
	if err != nil {
		return trade.Transaction{}, err
	}

	transaction.ID = meta.Key
	return transaction, nil
}

// List returns all documents in the collection.
func (r *repository) List(ctx context.Context) ([]trade.Transaction, error) {
	var err error
	results := []trade.Transaction{}

	query := fmt.Sprintf(
		`FOR entry IN %s
RETURN entry`, r.collectionName)
	cur, err := r.database.Query(ctx, query, map[string]interface{}{})
	if err != nil {
		return []trade.Transaction{}, err
	}
	defer cur.Close()

	for cur.HasMore() {
		var transaction trade.Transaction
		if _, err = cur.ReadDocument(ctx, &transaction); err != nil {
			return []trade.Transaction{}, err
		}
		results = append(results, transaction)
	}

	return results, nil
}

// Get returns transaction with given id. If no document is found returns NotFoundError.
func (r *repository) Get(ctx context.Context, id string) (trade.Transaction, error) {
	var err error
	col, err := r.database.Collection(ctx, r.collectionName)
	if err != nil {
		return trade.Transaction{}, err
	}

	var transaction trade.Transaction
	_, err = col.ReadDocument(ctx, id, &transaction)
	if err != nil {
		return trade.Transaction{}, err
	}

	return transaction, nil
}

// Update updates the document identified by id with values transaction and returns
// the new document. If no document with id is found returns NotFoundError.
func (r *repository) Update(ctx context.Context, id string, transaction trade.Transaction) (trade.Transaction, error) {
	var err error
	col, err := r.database.Collection(ctx, r.collectionName)
	if err != nil {
		return trade.Transaction{}, err
	}

	var result trade.Transaction
	_, err = col.UpdateDocument(arango.WithReturnNew(ctx, &result), id, transaction)
	if err != nil {
		return trade.Transaction{}, err
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
