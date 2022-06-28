package graphql

import (
	"github.com/ichaly/go-gql/transport"
	"github.com/ichaly/go-gql/types"
	"net/http"
)

type Server struct {
	transports []types.Transport
	exec       *types.Executor
}

func NewServer() *Server {
	srv := &Server{
		exec: &types.Executor{},
	}

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	return srv
}

func (s *Server) AddTransport(transport types.Transport) {
	s.transports = append(s.transports, transport)
}

func (s *Server) getTransport(r *http.Request) types.Transport {
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
			types.SendErrorf(w, http.StatusUnprocessableEntity, "internal system error")
		}
	}()

	if t := s.getTransport(r); t != nil {
		t.Do(w, r, s.exec)
		return
	}

	types.SendErrorf(w, http.StatusBadRequest, "transport not supported")
}
