package graphql

import (
	"context"
	"github.com/ichaly/go-gql/internal/executor"
)

type Schema struct {
}

func (s *Schema) Exec(
	ctx context.Context,
	queryString string,
	operationName string,
	variables map[string]interface{},
) *executor.GqlResult {
	return nil
}
