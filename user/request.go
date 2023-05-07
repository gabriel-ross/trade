package user

// TODO: handle request body
// TODO: handle request url query params
// TODO: separately define common url query params (sort, direction, limit, etc.) from route/resource-specific url params

// ?: Maybe something for unmarshaling URL elements into different structs?

type URLQueryParam struct {
	URLKeyName   string
	Required     bool
	DefaultValue interface{}
}
