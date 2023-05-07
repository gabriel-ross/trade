package trade

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
)

// sort query will be formatted ?sortKey=key1+dir,key2+dir
// TODO: converting all sort keys to a db query

// Default sort direction is ascending
// Default limit is 10
var (
	DEFAULT_LIMIT = 10
)

func Extract(r *http.Request) (URLQueryParams, error) {
	var err error
	queryVals := URLQueryParams{
		SortKeys: []SortKey{},
		Limit:    DEFAULT_LIMIT,
	}
	if sort := r.URL.Query().Get("sort"); sort != "" {
		keys := strings.Split(sort, ",")
		for _, key := range keys {
			vals := strings.Split(key, " ")
			newSortKey := SortKey{}
			if len(vals) > 0 {
				newSortKey.Field = vals[0]
			}
			if len(vals) > 1 && strings.ToLower(vals[1]) == "desc" {
				newSortKey.Direction = SORT_DESC
			} else {
				newSortKey.Direction = SORT_ASC
			}
			queryVals.SortKeys = append(queryVals.SortKeys, newSortKey)
		}
	}

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		queryVals.Limit, err = strconv.Atoi(limitStr)
		if err != nil {
			return URLQueryParams{}, err
		}
	}

	return queryVals, nil
}

type URLQueryParams struct {
	SortKeys []SortKey
	Limit    int
}

type SortKey struct {
	Field     string
	Direction SortDirection
}

func NewQueryMap(r *http.Request, queryParams map[string]bool) (map[string]string, error) {
	qm := map[string]string{}
	for queryParam, required := range queryParams {
		if qm[queryParam] = chi.URLParam(r, queryParam); qm[queryParam] == "" && required {
			return nil, fmt.Errorf("missing required parameter %s", queryParam)
		}
	}
	return qm, nil
}
