package graphql

import (
	"log"
	"testing"
)

type Todo struct {
	ID string
}

func TestBuild(t *testing.T) {
	builder := NewBuilder()
	q := builder.Query()
	q.Field("todos", func() Todo {
		return Todo{}
	})
	log.Printf("%v", builder.MustBuild())
}
