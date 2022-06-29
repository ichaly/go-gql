package graphql

import (
	"github.com/ichaly/go-gql/internal/executor"
	"net/http"
)

type Transport interface {
	Supports(r *http.Request) bool
	Do(w http.ResponseWriter, r *http.Request, exec *executor.Executor)
}

type Enum interface {
	String() string
	Values() []string
}

type Node interface {
	Alias() string
}
