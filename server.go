package graphql

import (
	"github.com/ichaly/go-gql/errors"
	"github.com/ichaly/go-gql/transport"
	"net/http"
	"time"
)

type (
	Transport interface {
		Supports(r *http.Request) bool
		Do(w http.ResponseWriter, r *http.Request, exec *Executor)
	}
	Server struct {
		transports []Transport
		exec       *Executor
	}
	RawParams struct {
		Query         string                 `json:"query"`
		OperationName string                 `json:"operationName"`
		Variables     map[string]interface{} `json:"variables"`
		Extensions    map[string]interface{} `json:"extensions"`
		Headers       http.Header            `json:"headers"`

		ReadTime TraceTiming `json:"-"`
	}
	TraceTiming struct {
		Start time.Time
		End   time.Time
	}
)

var Now = time.Now

func NewServer() *Server {
	srv := &Server{
		exec: &Executor{},
	}

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	return srv
}

func (s *Server) AddTransport(transport Transport) {
	s.transports = append(s.transports, transport)
}

func (s *Server) getTransport(r *http.Request) Transport {
	for _, t := range s.transports {
		if t.Supports(r) {
			return t
		}
	}
	return nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			errors.SendErrorf(w, http.StatusUnprocessableEntity, "internal system error")
		}
	}()

	if t := s.getTransport(r); t != nil {
		t.Do(w, r, s.exec)
		return
	}

	errors.SendErrorf(w, http.StatusBadRequest, "transport not supported")
}
