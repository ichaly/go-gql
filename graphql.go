package graphql

import (
	"context"
	"encoding/json"
	"github.com/ichaly/go-gql/errors"
)

type Schema struct {
}

type Response struct {
	Errors     []*errors.QueryError   `json:"errors,omitempty"`
	Data       json.RawMessage        `json:"data,omitempty"`
	Extensions map[string]interface{} `json:"extensions,omitempty"`
}

func (s *Schema) Exec(ctx context.Context, queryString string, operationName string, variables map[string]interface{}) *Response {
	return nil
}
