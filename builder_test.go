package graphql

import (
	"log"
	"testing"
)

type Todo struct {
	Finish bool
	Owner  *User
	Data   *interface{}
	//Api EntryPoint
	//createTime time.Time
}

type User struct {
	ID string
}

type EntryPoint interface {
}

func TestBuild(t *testing.T) {
	builder := NewBuilder()
	q := builder.Query()
	q.Field("todo", func() Todo {
		return Todo{}
	})
	q.Field("todos", func() []Todo {
		return []Todo{}
	})
	q.Field("ptr", func() *Todo {
		return &Todo{}
	})
	q.Field("ptrs", func() []*Todo {
		return []*Todo{}
	})
	log.Printf("%v", builder.MustBuild())
}
