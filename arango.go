package trade

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	arangodriver "github.com/arangodb/go-driver"
	arangohttp "github.com/arangodb/go-driver/http"
	"github.com/tidwall/gjson"
)

type Schema struct {
	DocumentCollections []string `json:"documentCollections"`
	EdgeCollections     []string `json:"edgeCollections"`
}

type ArangoClient struct {
	DriverConnection arangodriver.Connection
	DriverClient     arangodriver.Client
}

func NewArangoClient(addrs []string) (*ArangoClient, error) {
	conn, err := arangohttp.NewConnection(arangohttp.ConnectionConfig{
		Endpoints: addrs,
	})
	if err != nil {
		return nil, err
	}

	cl, err := arangodriver.NewClient(arangodriver.ClientConfig{
		Connection: conn,
	})
	if err != nil {
		return nil, err
	}
	return &ArangoClient{
		DriverClient: cl,
	}, nil
}

func (cl *ArangoClient) Database(ctx context.Context, name string, createOnNotExist bool, schemaPath string) (arangodriver.Database, error) {
	var err error
	dbClient, err := cl.DriverClient.Database(ctx, name)
	if err != nil {
		if arangodriver.IsNotFoundGeneral(err) && createOnNotExist {
			// Create database with name
			dbClient, err := cl.DriverClient.CreateDatabase(ctx, name, nil)
			if err != nil {
				return nil, err
			}

			schemaF, err := os.Open(schemaPath)
			if err != nil {
				return nil, err
			}

			schemaRaw, err := io.ReadAll(schemaF)
			if err != nil {
				return nil, err
			}

			var schema Schema
			err = json.Unmarshal(schemaRaw, &schema)
			if err != nil {
				return nil, err
			}

			// Create document collections
			for _, e := range schema.DocumentCollections {
				_, err = dbClient.CreateCollection(ctx, gjson.Get(e, "collectionName").Str, nil)
				if err != nil {
					return nil, err
				}
			}

			// Create edge collections
			for _, e := range schema.EdgeCollections {
				_, err = dbClient.CreateCollection(ctx, gjson.Get(e, "collectionName").Str, &arangodriver.CreateCollectionOptions{Type: arangodriver.CollectionTypeEdge})
				if err != nil {
					return nil, err
				}
			}

		} else {
			return nil, err
		}
	}
	return dbClient, nil
}

type ArangoRepository[T any] struct {
	database       arangodriver.Database
	collectionName string
}

func NewArangoRepository[T any](db arangodriver.Database, collectionName string) *ArangoRepository[T] {
	return &ArangoRepository[T]{
		database:       db,
		collectionName: collectionName,
	}
}

func (r *ArangoRepository[T]) Create(ctx context.Context, data T) (string, T, error) {
	var err error
	var t T
	col, err := r.database.Collection(ctx, r.collectionName)
	if err != nil {
		return "", t, err
	}

	meta, err := col.CreateDocument(ctx, data)
	if err != nil {
		return "", t, err
	}

	return meta.Key, data, nil
}

func (r *ArangoRepository[T]) Query(ctx context.Context, query string) ([]T, error) {
	var err error
	results := []T{}

	cur, err := r.database.Query(ctx, query, map[string]interface{}{})
	if err != nil {
		return nil, err
	}
	defer cur.Close()

	for cur.HasMore() {
		var data T
		if _, err = cur.ReadDocument(ctx, &data); err != nil {
			return nil, err
		}
		results = append(results, data)
	}
	return results, nil
}

func (r *ArangoRepository[T]) Get(ctx context.Context, id string) (T, error) {
	var err error
	var t T

	col, err := r.database.Collection(ctx, r.collectionName)
	if err != nil {
		return t, err
	}

	var result T
	_, err = col.ReadDocument(ctx, id, &result)
	if err != nil {
		return t, err
	}

	return result, nil
}

// Update updates the document identified by id with values user and returns
// the new document. If no document with id is found returns NotFoundError.
func (r *ArangoRepository[T]) Update(ctx context.Context, id string, data T) (T, error) {
	var err error
	var t T
	col, err := r.database.Collection(ctx, r.collectionName)
	if err != nil {
		return t, err
	}

	var result T
	_, err = col.UpdateDocument(arangodriver.WithReturnNew(ctx, &result), id, data)
	if err != nil {
		return t, err
	}

	return result, nil
}

// Delete deletes the document with given id from the collection. If no match
// is found returns NotFoundError.
func (r *ArangoRepository[T]) Delete(ctx context.Context, id string) error {
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

// sort query will be formatted ?sort=key1+dir,key2+dir
// Operator, value separated by + which decodes to a space. If there is no operator default to equality
func BuildFilterQueryFromURLParams(aqb ArangoQueryBuilder, r *http.Request, queryParams []string, paginate Paginate) (ArangoQueryBuilder, error) {
	if len(queryParams) < 1 {
		return aqb, nil
	}

	i := 0
	fqb := NewFilterQueryBuilder(aqb)
	for idx, param := range queryParams {
		if val := r.URL.Query().Get(param); val != "" {
			fqb = aqb.Filter(FilterKeyFromURLElement(param, val))
			i = idx
			break
		}
	}

	if i+1 < len(queryParams) {
		for _, param := range queryParams[i+1:] {
			if val := r.URL.Query().Get(param); val != "" {
				if r.URL.Query().Get("inclusive") == "true" {
					fqb = fqb.Or(FilterKeyFromURLElement(param, val))
				} else {
					fqb = fqb.And(FilterKeyFromURLElement(param, val))
				}
			}
		}
	}

	aqb = fqb.Sort(paginate.SortFields...)
	aqb = aqb.Limit(paginate.Limit)

	// extract sort and limit values. Move ExtractURLQueryParams code from request to here
	return fqb.Sort(paginate.SortFields...).Limit(paginate.Limit).Done(), nil
}

type SortDirection string

type FilterOperator string

var (
	SORT_ASC     = SortDirection("ASC")
	SORT_DESC    = SortDirection("DESC")
	Eq           = FilterOperator("==")
	Neq          = FilterOperator("!=")
	Gt           = FilterOperator(">")
	Lt           = FilterOperator("<")
	Geq          = FilterOperator(">=")
	Leq          = FilterOperator("<=")
	OPERATOR_MAP = map[string]FilterOperator{
		"eq":  Eq,
		"neq": Neq,
		"gt":  Gt,
		"lt":  Lt,
		"geq": Geq,
		"leq": Leq,
	}
)

type FilterKey struct {
	FieldName string
	Operator  FilterOperator
	Value     interface{}
}

func NewFilterKey(fieldName string, op FilterOperator, val interface{}) FilterKey {
	return FilterKey{
		FieldName: fieldName,
		Operator:  op,
		Value:     val,
	}
}

func (fk FilterKey) FormatValue() interface{} {
	switch v := fk.Value.(type) {
	case string:
		return fmt.Sprintf("\"%s\"", v)
	default:
		return fk.Value
	}
}

func FilterKeyFromURLElement(key, val string) FilterKey {
	eles := strings.Split(val, " ")
	if len(eles) < 2 {
		return FilterKey{
			FieldName: key,
			Operator:  Eq,
			Value:     eles[0],
		}
	} else {
		return FilterKey{
			FieldName: key,
			Operator:  OPERATOR_MAP[eles[0]],
			Value:     eles[1],
		}
	}
}

type ArangoQueryBuilder struct {
	QueryString *strings.Builder
	loopVar     string
}

func NewArangoQueryBuilder(collectionName string) ArangoQueryBuilder {
	aqb := ArangoQueryBuilder{
		QueryString: &strings.Builder{},
		loopVar:     "x",
	}
	aqb.QueryString.WriteString(fmt.Sprintf("FOR %s IN %s", aqb.loopVar, collectionName))
	return aqb
}

func (aqb ArangoQueryBuilder) Sort(sortFields ...SortField) ArangoQueryBuilder {
	if len(sortFields) > 0 {
		aqb.QueryString.WriteString(fmt.Sprintf("\n\tSORT %s.%s %s", aqb.loopVar, sortFields[0].Field, sortFields[0].Direction))
	}
	if len(sortFields) > 1 {
		for _, sortField := range sortFields {
			aqb.QueryString.WriteString(fmt.Sprintf(", %s.%s %s", aqb.loopVar, sortField.Field, sortField.Direction))
		}
	}
	return aqb
}

func (aqb ArangoQueryBuilder) Limit(limit int) ArangoQueryBuilder {
	aqb.QueryString.WriteString(fmt.Sprintf("\n\tLIMIT %d", limit))
	return aqb
}

func (aqb ArangoQueryBuilder) Paginate(p Paginate) ArangoQueryBuilder {
	return aqb.Sort(p.SortFields...).Limit(p.Limit)
}

func (aqb ArangoQueryBuilder) Filter(filter FilterKey) FilterQueryBuilder {
	aqb.QueryString.WriteString(fmt.Sprintf("\n\tFILTER %s.%s %s %v", aqb.loopVar, filter.FieldName, filter.Operator, filter.FormatValue()))
	return FilterQueryBuilder{aqb}
}

func (aqb ArangoQueryBuilder) Done() ArangoQueryBuilder {
	aqb.QueryString.WriteString(fmt.Sprintf("\n\tRETURN %s", aqb.loopVar))
	return aqb
}

func (aqb ArangoQueryBuilder) String() string {
	return aqb.QueryString.String()
}

type FilterQueryBuilder struct {
	ArangoQueryBuilder
}

func NewFilterQueryBuilder(aqb ArangoQueryBuilder) FilterQueryBuilder {
	return FilterQueryBuilder{aqb}
}

func (fqb FilterQueryBuilder) And(filter FilterKey) FilterQueryBuilder {
	fqb.QueryString.WriteString(fmt.Sprintf(" && %s.%s %s %v", fqb.loopVar, filter.FieldName, filter.Operator, filter.FormatValue()))
	return fqb
}

func (fqb FilterQueryBuilder) Or(filter FilterKey) FilterQueryBuilder {
	fqb.QueryString.WriteString(fmt.Sprintf(" || %s.%s %s %v", fqb.loopVar, filter.FieldName, filter.Operator, filter.FormatValue()))
	return fqb
}

func (fqb FilterQueryBuilder) Not(filter FilterKey) FilterQueryBuilder {
	fqb.QueryString.WriteString(fmt.Sprintf(" NOT %s.%s %s %v", fqb.loopVar, filter.FieldName, filter.Operator, filter.FormatValue()))
	return fqb
}
