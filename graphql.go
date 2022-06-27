package graphql

import (
	"context"
)

type Schema struct {
}

func (s *Schema) Exec(ctx context.Context, queryString string, operationName string, variables map[string]interface{}) *Response {
	return nil
}
