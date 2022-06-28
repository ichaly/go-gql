package graphql

import "github.com/ichaly/go-gql/internal/introspection"

type Executor struct {
	introspection.Schema
}
