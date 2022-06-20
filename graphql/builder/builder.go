package builder

import (
	"github.com/ichaly/go-gql/graphql"
	"reflect"
)

type SchemaBuilder struct {
	Name    string
	objects map[string]*Object
}

type BuildOption interface {
	apply(*SchemaBuilder)
}

func NewSchemaBuilder(options ...BuildOption) *SchemaBuilder {
	builder := &SchemaBuilder{
		objects: make(map[string]*Object),
	}
	for _, o := range options {
		o.apply(builder)
	}
	return builder
}

type query struct{}

func (my *SchemaBuilder) Query() *Object {
	return my.Object("Query", query{})
}

type mutation struct{}

func (my *SchemaBuilder) Mutation() *Object {
	return my.Object("Mutation", mutation{})
}

type ObjectOption interface {
	apply(*SchemaBuilder, *Object)
}

func (my *SchemaBuilder) Object(name string, typ interface{}, options ...ObjectOption) *Object {
	if object, ok := my.objects[name]; ok {
		if reflect.TypeOf(object.Type) != reflect.TypeOf(typ) {
			panic("re-registered object with different type")
		}
		return object
	}
	object := &Object{
		Name:        name,
		Type:        typ,
		ServiceName: my.Name,
	}
	my.objects[name] = object

	for _, o := range options {
		o.apply(my, object)
	}

	return object
}

func (my *SchemaBuilder) Build() (*graphql.Schema, error) {
	return &graphql.Schema{}, nil
}

func (my *SchemaBuilder) MustBuild() *graphql.Schema {
	schema, err := my.Build()
	if err != nil {
		panic(err)
	}
	return schema
}
