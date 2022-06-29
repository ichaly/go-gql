package graphql

import (
	"github.com/ichaly/go-gql/internal/executor"
	transport2 "github.com/ichaly/go-gql/internal/transport"
	"github.com/ichaly/go-gql/internal/util"
	"net/http"
)

type Server struct {
	transports []Transport
	exec       *executor.Executor
}

func NewServer() *Server {
	srv := &Server{
		exec: &executor.Executor{},
	}

	srv.AddTransport(transport2.Options{})
	srv.AddTransport(transport2.GET{})
	srv.AddTransport(transport2.POST{})

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
			util.SendErrorf(w, http.StatusUnprocessableEntity, "internal system error")
		}
	}()

	if t := s.getTransport(r); t != nil {
		t.Do(w, r, s.exec)
		return
	}

	util.SendErrorf(w, http.StatusBadRequest, "transport not supported")
}
