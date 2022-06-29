package graphql

import (
	"github.com/ichaly/go-gql/internal/introspection"
	"reflect"
)

type Object struct {
	key         string // Optional, defaults 'ID'.
	Name        string // Optional, defaults to Type's name.
	Description string
	Type        interface{}
	Methods     map[string]interface{}
}

func (s *Object) Field(name string, method interface{}) {
	if s.Methods == nil {
		s.Methods = make(map[string]interface{})
	}

	if _, ok := s.Methods[name]; ok {
		panic("duplicate Method")
	}
	s.Methods[name] = method
}

type Builder struct {
	Name    string
	objects map[reflect.Type]*Object
}

type BuildOption interface {
	apply(*Builder)
}

func NewBuilder(options ...BuildOption) *Builder {
	builder := &Builder{
		objects: make(map[reflect.Type]*Object),
	}
	for _, o := range options {
		o.apply(builder)
	}
	return builder
}

type query struct{}

func (my *Builder) Query() *Object {
	return my.Object("Query", query{})
}

type mutation struct{}

func (my *Builder) Mutation() *Object {
	return my.Object("Mutation", mutation{})
}

type ObjectOption interface {
	apply(*Builder, *Object)
}

func (my *Builder) Object(name string, typ interface{}, options ...ObjectOption) *Object {
	if object, ok := my.objects[reflect.TypeOf(typ)]; ok {
		if reflect.TypeOf(object.Type) != reflect.TypeOf(typ) {
			panic("re-registered object with different type")
		}
		return object
	}
	object := &Object{
		Name: name,
		Type: typ,
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
