package trade

import (
	"context"
	"fmt"
	"strings"

	arangodriver "github.com/arangodb/go-driver"
	arangohttp "github.com/arangodb/go-driver/http"
)

type Client struct {
	DriverConnection arangodriver.Connection
	DriverClient     arangodriver.Client
}

func NewArangoClient(addrs []string) (*Client, error) {
	conn, err := arangohttp.NewConnection(arangohttp.ConnectionConfig{
		Endpoints: addrs,
	})
	if err != nil {
		return &Client{}, err
	}

	cl, err := arangodriver.NewClient(arangodriver.ClientConfig{
		Connection: conn,
	})
	if err != nil {
		return &Client{}, err
	}
	return &Client{DriverClient: cl}, nil
}

func (cl *Client) Database(ctx context.Context, name string, createOnNotExist bool, collections []string) (arangodriver.Database, error) {
	var err error
	dbClient, err := cl.DriverClient.Database(ctx, name)
	if err != nil {
		if arangodriver.IsNotFoundGeneral(err) && createOnNotExist {
			// Create database with name
			dbClient, err := cl.DriverClient.CreateDatabase(ctx, name, nil)
			if err != nil {
				return nil, err
			}
			// Create collections
			for _, colName := range collections {
				_, err = dbClient.CreateCollection(ctx, colName, nil)
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

type SortDirection string

type FilterOperator string

var (
	SORT_ASC  = SortDirection("ASC")
	SORT_DESC = SortDirection("DESC")
	Eq        = FilterOperator("==")
	Neq       = FilterOperator("!=")
	Gt        = FilterOperator(">")
	Lt        = FilterOperator("<")
	Geq       = FilterOperator(">=")
	Leq       = FilterOperator("<=")
)

type SortField struct {
	FieldName string
	Direction SortDirection
}

type FilterKey struct {
	FieldName string
	Operator  FilterOperator
	Value     interface{}
}

type ArangoQueryBuilder struct {
	QueryString *strings.Builder
	loopVar     string
}

func NewArangoQueryBuilder(collectionName string) *ArangoQueryBuilder {
	aqb := &ArangoQueryBuilder{
		QueryString: &strings.Builder{},
		loopVar:     "x",
	}
	aqb.QueryString.WriteString(fmt.Sprintf("FOR %s IN %s", aqb.loopVar, collectionName))
	return aqb
}

func (aqb *ArangoQueryBuilder) Sort(sortFields ...SortField) *ArangoQueryBuilder {
	aqb.QueryString.WriteString("\n\tSORT")
	for _, sortField := range sortFields {
		aqb.QueryString.WriteString(fmt.Sprintf(", %s.%s %s", aqb.loopVar, sortField.FieldName, sortField.Direction))
	}
	return aqb
}

func (aqb *ArangoQueryBuilder) Limit(limit int) *ArangoQueryBuilder {
	aqb.QueryString.WriteString(fmt.Sprintf("\n\tLIMIT %d", limit))
	return aqb
}

func (aqb *ArangoQueryBuilder) Filter() *FilterQueryBuilder {
	aqb.QueryString.WriteString("\n\tFILTER")
	return &FilterQueryBuilder{aqb}
}

func (aqb *ArangoQueryBuilder) String() string {
	aqb.QueryString.WriteString(fmt.Sprintf("\n\tRETURN %s", aqb.loopVar))
	return aqb.QueryString.String()
}

type FilterQueryBuilder struct {
	*ArangoQueryBuilder
}

func (fqb *FilterQueryBuilder) On(filter FilterKey) *FilterQueryBuilder {
	fqb.QueryString.WriteString(fmt.Sprintf(" %s.%s %s %v", fqb.loopVar, filter.FieldName, filter.Operator, filter.Value))
	return fqb
}

func (fqb *FilterQueryBuilder) And(filter FilterKey) *FilterQueryBuilder {
	fqb.QueryString.WriteString(fmt.Sprintf(" && %s.%s %s %v", fqb.loopVar, filter.FieldName, filter.Operator, filter.Value))
	return fqb
}

func (fqb *FilterQueryBuilder) Or(filter FilterKey) *FilterQueryBuilder {
	fqb.QueryString.WriteString(fmt.Sprintf(" || %s.%s %s %v", fqb.loopVar, filter.FieldName, filter.Operator, filter.Value))
	return fqb
}

func (fqb *FilterQueryBuilder) Not(filter FilterKey) *FilterQueryBuilder {
	fqb.QueryString.WriteString(fmt.Sprintf(" NOT %s.%s %s %v", fqb.loopVar, filter.FieldName, filter.Operator, filter.Value))
	return fqb
}
