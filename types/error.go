package types

import (
	"encoding/json"
	"fmt"
	"github.com/ichaly/go-gql"
	"net/http"
)

type GqlError struct {
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
	SendError(w, code, &GqlError{Message: fmt.Sprintf(format, args...)})
}

func SendError(w http.ResponseWriter, code int, errors ...*GqlError) {
	w.WriteHeader(code)
	b, err := json.Marshal(&graphql.GqlResponse{Errors: errors})
	if err != nil {
		panic(err)
	}
	_, _ = w.Write(b)
}
