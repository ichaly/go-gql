package types

import (
	"encoding/json"
	"fmt"
	"github.com/ichaly/go-gql/internal/introspection"
	"net/http"
	"time"
)

type (
	Executor struct {
		introspection.Schema
	}
	Transport interface {
		Supports(r *http.Request) bool
		Do(w http.ResponseWriter, r *http.Request, exec *Executor)
	}
	TraceTiming struct {
		Start time.Time
		End   time.Time
	}
	GqlRequest struct {
		Query         string                 `json:"query"`
		OperationName string                 `json:"operationName"`
		Variables     map[string]interface{} `json:"variables"`
		Extensions    map[string]interface{} `json:"extensions"`
		Headers       http.Header            `json:"headers"`

		ReadTime TraceTiming `json:"-"`
	}
	GqlResponse struct {
		Errors     []*GqlError            `json:"errors,omitempty"`
		Data       json.RawMessage        `json:"data,omitempty"`
		Extensions map[string]interface{} `json:"extensions,omitempty"`
	}
	GqlError struct {
		Err           error                  `json:"-"` // Err holds underlying if available
		Message       string                 `json:"message"`
		Locations     []Location             `json:"locations,omitempty"`
		Path          []interface{}          `json:"path,omitempty"`
		Rule          string                 `json:"-"`
		ResolverError error                  `json:"-"`
		Extensions    map[string]interface{} `json:"extensions,omitempty"`
	}
	Location struct {
		Line   int `json:"line"`
		Column int `json:"column"`
	}
)

func SendErrorf(w http.ResponseWriter, code int, format string, args ...interface{}) {
	SendError(w, code, &GqlError{Message: fmt.Sprintf(format, args...)})
}

func SendError(w http.ResponseWriter, code int, errors ...*GqlError) {
	w.WriteHeader(code)
	b, err := json.Marshal(&GqlResponse{Errors: errors})
	if err != nil {
		panic(err)
	}
	_, _ = w.Write(b)
}
