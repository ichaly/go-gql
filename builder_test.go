package graphql

import (
	"log"
	"testing"
)

type Todo struct {
	ID     *string
	Finish bool
}

func TestBuild(t *testing.T) {
	builder := NewBuilder()
	q := builder.Query()
	q.Field("ptrs", func() []*Todo {
		return []*Todo{}
	})
	q.Field("todo", func() Todo {
		return Todo{}
	})
	q.Field("ptr", func() *Todo {
		return &Todo{}
	})
	q.Field("todos", func() []Todo {
		return []Todo{}
	})
	log.Printf("%v", builder.MustBuild())
}
