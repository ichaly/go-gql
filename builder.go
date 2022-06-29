package graphql

import (
	"github.com/vektah/gqlparser/v2"
	"github.com/vektah/gqlparser/v2/ast"
	"log"
	"reflect"
	"strings"
)

type Object struct {
	key         string // Optional, defaults 'ID'.
	Name        string // Optional, defaults to Type's name.
	Description string
	Type        interface{}
	Resolvers   map[string]*Object
}

type FieldOption interface {
	apply(*Object)
}

func (s *Object) Field(name string, resolver interface{}, options ...FieldOption) {
	if s.Resolvers == nil {
		s.Resolvers = make(map[string]*Object)
	}
	if _, ok := s.Resolvers[name]; ok {
		panic("duplicate Method")
	}
	m := &Object{Type: resolver}
	for _, o := range options {
		o.apply(m)
	}
	s.Resolvers[name] = m
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

func (my *Builder) MustBuild() *ast.Schema {
	schema, err := my.Build()
	if err != nil {
		panic(err)
	}
	return schema
}

func (my *Builder) Build() (*ast.Schema, error) {
	sb := &strings.Builder{} // where the (text) schema is generated
	sb.Grow(256)             // Even simple schemas are at least this big

	sb.WriteString(getSchema(my.Query()))

	log.Print(sb.String())

	return gqlparser.LoadSchema(&ast.Source{
		Name:  "schema",
		Input: sb.String(),
	})
}

func getDescription(desc string) string {
	sb := &strings.Builder{}
	if len(desc) > 0 {
		sb.WriteString(`"""`)
		sb.WriteString(desc)
		sb.WriteString(`"""`)
		sb.WriteRune('\n')
	}
	return sb.String()
}

func getType(t reflect.Type) (result reflect.Type, nullable bool, iterable bool) {
	switch t.Kind() {
	case reflect.Struct:
		result = t
		return
	case reflect.Ptr:
		result = t.Elem()
		nullable = true
		return
	case reflect.Map, reflect.Slice, reflect.Array:
		result = t.Elem()
		iterable = true
		return
	case reflect.Func:
		if t.NumOut() == 0 {
			panic("Resolver func must have at least one return value")
		}
		return getType(t.Out(0))
	}
	return
}

func getSchema(o *Object) string {
	sb := &strings.Builder{}
	sb.WriteString(getDescription(o.Description))

	sb.WriteString("type ")
	sb.WriteString(o.Name)
	sb.WriteString(" { ")
	sb.WriteRune('\n')
	for k, v := range o.Resolvers {
		sb.WriteString(getDescription(v.Description))
		sb.WriteString(k)
		sb.WriteString(": ")
		t, n, i := getType(reflect.TypeOf(v.Type))
		if t == nil {
			continue
		}
		if i {
			sb.WriteString("[")
		}
		sb.WriteString(t.Name())
		if n {
			sb.WriteString("!")
		}
		if i {
			sb.WriteString("]!")
		}
		sb.WriteRune('\n')
	}
	sb.WriteString("}\n\n")

	//sb.WriteString("type Todo {id: ID!}")
	return sb.String()
}
