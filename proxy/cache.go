package proxy

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"
)

var ErrNotFound = errors.New("error request not found in cache")
var ErrStaleCachedState = errors.New("error cached state is stale")

// Response wraps an http.Response with a Timestamp of when it was cached.
type Response struct {
	Timestamp time.Time
	response  http.Response
}

// key returns a content-based key for an http request.
func key(r *http.Request) (string, error) {
	var key strings.Builder
	var err error
	if _, err = writeStrings(&key, r.Method, r.RequestURI); err != nil {
		return "", err
	}

	// Write headers to key
	for name, values := range r.Header {
		_, err = key.WriteString(name)
		if err != nil {
			return "", err
		}
		for _, value := range values {
			_, err = key.WriteString(value)
			if err != nil {
				return "", err
			}
		}
	}

	// Write body to key and restore request body
	buf, err := io.ReadAll(r.Body)
	if err != nil {
		return "", err
	}
	r.Body = io.NopCloser(bytes.NewBuffer(buf))
	_, err = key.Write(buf)
	if err != nil {
		return "", err
	}

	return key.String(), nil
}

// writeStrings is a helper function for writing multiple strings to
// a strings.Builder
func writeStrings(builder *strings.Builder, strings ...string) (int, error) {
	totalWritten := 0
	for _, s := range strings {
		written, err := builder.WriteString(s)
		if err != nil {
			return 0, err
		}
		totalWritten += written
	}
	return totalWritten, nil
}

type cache struct {
	contents map[string]Response
	timeout  time.Duration
}

func NewCache(timeout time.Duration) *cache {
	return &cache{
		contents: map[string]Response{},
		timeout:  timeout,
	}
}

func (c *cache) Fetch(r *http.Request) (Response, error) {
	var err error
	reqKey, err := key(r)
	if err != nil {
		return Response{}, err
	}
	cachedResp, exists := c.contents[reqKey]
	switch true {
	case !exists:
		return Response{}, ErrNotFound
	case exists && time.Since(cachedResp.Timestamp) >= c.timeout:
		return Response{}, ErrStaleCachedState
	}
	return cachedResp, nil
}

func (c *cache) Upsert(req *http.Request, resp http.Response) error {
	var err error
	reqKey, err := key(req)
	if err != nil {
		return err
	}

	c.contents[reqKey] = Response{
		Timestamp: time.Now(),
		response:  resp,
	}

	return nil
}
