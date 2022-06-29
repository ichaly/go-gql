package types

import (
	"github.com/ichaly/go-gql/executor"
	"net/http"
)

type Transport interface {
	Supports(r *http.Request) bool
	Do(w http.ResponseWriter, r *http.Request, exec *executor.Executor)
}
