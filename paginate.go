package trade

import (
	"net/http"
	"strconv"
	"strings"
)

var (
	DEFAULT_LIMIT = 1000
)

type Paginate struct {
	SortFields []SortField
	Limit      int
	// something for pagination
}

// NewPaginate creates a new paginate from pagination query values in an http
// request.
//
// NewPaginate expects sort quries in the format ?sort=key1+sortDirection,key2+sortDirection,etc.
func NewPaginate(r *http.Request) Paginate {
	p := Paginate{
		SortFields: []SortField{},
		Limit:      DEFAULT_LIMIT,
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

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		p.Limit, _ = strconv.Atoi(limitStr)
	}

	return p
}

type SortField struct {
	Field     string
	Direction SortDirection
}
