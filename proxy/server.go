package proxy

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type Config struct {
	NAME           string
	PORT           string `env:"PORT" default:"8081" required:"false"`
	SERVER_ADDRESS string `env:"SERVER_ADDRESS" default:"localhost:8080" required:"false"`
	CACHE_TIMEOUT  time.Duration
}

type Cache interface {

	// Fetch fetches the cached response for a request. If no cached data is
	// found returns ErrNotFound. If cached response is found and stale
	// returns ErrStaleCachedState.
	Fetch(r *http.Request) (Response, error)

	// Upsert caches an http response or updates stale data.
	Upsert(req *http.Request, resp http.Response) error
}

type Server struct {
	cnf    Config
	client http.Client
	cache  *cache
}

func New(cnf Config) *Server {
	return &Server{
		cnf:    cnf,
		cache:  NewCache(cnf.CACHE_TIMEOUT),
		client: *http.DefaultClient,
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var err error
	cachedResp, err := s.cache.Fetch(r)
	if errors.Is(err, ErrNotFound) || errors.Is(err, ErrStaleCachedState) {
		fwdReq, err := http.NewRequest(r.Method, s.cnf.SERVER_ADDRESS, r.Body)
		if err != nil {
			log.Println(err)
			return
		}
		// ?: Should I forward all the headers too?
		fwdReq.Header = r.Header

		// Forward request to server
		resp, err := s.client.Do(fwdReq)
		if err != nil {
			log.Println(err)
			return
		}

		// Cache response
		err = s.cache.Upsert(r, *resp)
		if err != nil {
			log.Println(err)
			return
		}

		// Forward response to client
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
			return
		}
		for name, values := range resp.Header {
			for _, value := range values {
				w.Header().Add(name, value)
			}
		}
		w.WriteHeader(resp.StatusCode)
		w.Write(body)
		return

	} else if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}

	body, err := io.ReadAll(cachedResp.response.Body)
	if err != nil {
		log.Println(err)
		return
	}
	for name, values := range cachedResp.response.Header {
		for _, value := range values {
			w.Header().Add(name, value)
		}
	}
	w.WriteHeader(cachedResp.response.StatusCode)
	w.Write(body)
}

func (s *Server) Run() error {
	log.Printf("%s server running on port %s", s.cnf.NAME, s.cnf.PORT)
	return http.ListenAndServe(":"+s.cnf.PORT, s)
}

func Ping() http.HandlerFunc {
	resp := fmt.Sprintf("Ping received. Server is healthy")
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(resp))
	}
}
