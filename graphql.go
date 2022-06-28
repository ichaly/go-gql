package graphql

import (
	"context"
	"github.com/ichaly/go-gql/types"
)

type Schema struct {
}

func (s *Schema) Exec(
	ctx context.Context,
	queryString string,
	operationName string,
	variables map[string]interface{},
) *types.GqlResponse {
	return nil
}
