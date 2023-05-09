package trade

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
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

// sort query will be formatted ?sort=key1+dir,key2+dir
// Operator, value separated by + which decodes to a space. If there is no operator default to equality
func BuildFilterQueryFromURLParams(aqb ArangoQueryBuilder, r *http.Request, queryParams []string) (ArangoQueryBuilder, error) {
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

	sortFields := []SortField{}
	if sort := r.URL.Query().Get("sort"); sort != "" {
		keys := strings.Split(sort, ",")
		for _, key := range keys {
			vals := strings.Split(key, " ")
			sf := SortField{}
			if len(vals) > 0 {
				sf.Field = vals[0]
			}
			if len(vals) > 1 && strings.ToLower(vals[1]) == "desc" {
				sf.Direction = SORT_DESC
			} else {
				sf.Direction = SORT_ASC
			}
			sortFields = append(sortFields, sf)
		}
	}
	aqb = fqb.Sort(sortFields...)

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			return ArangoQueryBuilder{}, err
		}
		aqb = fqb.Limit(limit)
	}

	// extract sort and limit values. Move ExtractURLQueryParams code from request to here
	return aqb.Done(), nil
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

type SortField struct {
	Field     string
	Direction SortDirection
}

type FilterKey struct {
	FieldName string
	Operator  FilterOperator
	Value     interface{}
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

func (aqb ArangoQueryBuilder) Filter(filter FilterKey) FilterQueryBuilder {
	aqb.QueryString.WriteString(fmt.Sprintf("\n\tFILTER %s.%s %s \"%v\"", aqb.loopVar, filter.FieldName, filter.Operator, filter.Value))
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
	fqb.QueryString.WriteString(fmt.Sprintf(" && %s.%s %s \"%v\"", fqb.loopVar, filter.FieldName, filter.Operator, filter.Value))
	return fqb
}

func (fqb FilterQueryBuilder) Or(filter FilterKey) FilterQueryBuilder {
	fqb.QueryString.WriteString(fmt.Sprintf(" || %s.%s %s \"%v\"", fqb.loopVar, filter.FieldName, filter.Operator, filter.Value))
	return fqb
}

func (fqb FilterQueryBuilder) Not(filter FilterKey) FilterQueryBuilder {
	fqb.QueryString.WriteString(fmt.Sprintf(" NOT %s.%s %s \"%v\"", fqb.loopVar, filter.FieldName, filter.Operator, filter.Value))
	return fqb
}
