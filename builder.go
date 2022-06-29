package graphql

import (
	"github.com/ichaly/go-gql/internal/introspection"
	"github.com/ichaly/go-gql/types"
	"reflect"
)

type Builder struct {
	Name    string
	objects map[reflect.Type]*types.Object
}

type BuildOption interface {
	apply(*Builder)
}

func NewBuilder(options ...BuildOption) *Builder {
	builder := &Builder{
		objects: make(map[reflect.Type]*types.Object),
	}
	for _, o := range options {
		o.apply(builder)
	}
	return builder
}

type query struct{}

func (my *Builder) Query() *types.Object {
	return my.Object("Query", query{})
}

type mutation struct{}

func (my *Builder) Mutation() *types.Object {
	return my.Object("Mutation", mutation{})
}

type ObjectOption interface {
	apply(*Builder, *types.Object)
}

func (my *Builder) Object(name string, typ interface{}, options ...ObjectOption) *types.Object {
	if object, ok := my.objects[reflect.TypeOf(typ)]; ok {
		if reflect.TypeOf(object.Type) != reflect.TypeOf(typ) {
			panic("re-registered object with different type")
		}
		return object
	}
	object := &types.Object{
		Name:        name,
		Type:        typ,
		ServiceName: my.Name,
	}
	my.objects[reflect.TypeOf(typ)] = object

	for _, o := range options {
		o.apply(my, object)
	}

	return object
}

func (my *Builder) MustBuild() *introspection.Schema {
	schema, err := my.Build()
	if err != nil {
		panic(err)
	}
	return schema
}

func (my *Builder) Build() (*introspection.Schema, error) {
	my.Object("Query", query{})
	my.Object("Mutation", mutation{})
	return &introspection.Schema{
		QueryType:    nil,
		MutationType: nil,
	}, nil
}
