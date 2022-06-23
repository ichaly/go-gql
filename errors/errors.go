package errors

import (
	"encoding/json"
	"fmt"
	graphql "github.com/ichaly/go-gql"
	"net/http"
)

type QueryError struct {
	Err           error                  `json:"-"` // Err holds underlying if available
	Message       string                 `json:"message"`
	Locations     []Location             `json:"locations,omitempty"`
	Path          []interface{}          `json:"path,omitempty"`
	Rule          string                 `json:"-"`
	ResolverError error                  `json:"-"`
	Extensions    map[string]interface{} `json:"extensions,omitempty"`
}

type Location struct {
	Line   int `json:"line"`
	Column int `json:"column"`
}

func SendErrorf(w http.ResponseWriter, code int, format string, args ...interface{}) {
	SendError(w, code, &QueryError{Message: fmt.Sprintf(format, args...)})
}

func SendError(w http.ResponseWriter, code int, errors ...*QueryError) {
	w.WriteHeader(code)
	b, err := json.Marshal(&graphql.Response{Errors: errors})
	if err != nil {
		panic(err)
	}
	_, _ = w.Write(b)
}
