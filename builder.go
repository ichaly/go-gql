package graphql

import (
	"fmt"
	"github.com/vektah/gqlparser/v2"
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/gqlerror"
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
	return my.Object(query{}, WithName("Query"))
}

type mutation struct{}

func (my *Builder) Mutation() *Object {
	return my.Object(mutation{}, WithName("Mutation"))
}

type objectOption func(*Builder, *Object)

func WithName(name string) objectOption {
	return func(b *Builder, o *Object) {
		o.Name = name
	}
}

func (my *Builder) Object(typ interface{}, options ...objectOption) *Object {
	if object, ok := my.objects[reflect.TypeOf(typ)]; ok {
		if reflect.TypeOf(object.Type) != reflect.TypeOf(typ) {
			panic("re-registered object with different type")
		}
		return object
	}
	var name string
	if v, ok := typ.(reflect.Type); ok {
		name = v.Name()
	} else {
		name = reflect.TypeOf(typ).Name()
	}
	object := &Object{
		Name: name,
		Type: typ,
	}
	my.objects[reflect.TypeOf(typ)] = object

	for _, o := range options {
		o(my, object)
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

func (my *Builder) Build() (*ast.Schema, *gqlerror.Error) {
	sb := &strings.Builder{}

	for _, o := range my.objects {
		sb.WriteString(my.getSchema(o))
	}

	fmt.Println(sb.String())
	return gqlparser.LoadSchema(&ast.Source{
		Name:  "schema",
		Input: sb.String(),
	})
}

func (my *Builder) getDescription(desc string) string {
	sb := &strings.Builder{}
	if len(desc) > 0 {
		sb.WriteString(`"""`)
		sb.WriteString(desc)
		sb.WriteString(`"""`)
		sb.WriteRune('\n')
	}
	return sb.String()
}

func (my *Builder) getType(t reflect.Type) (result reflect.Type, nullable bool, iterable bool) {
	for {
		switch t.Kind() {
		case reflect.Struct:
			result = t
			my.Object(result)
			return
		case reflect.Ptr:
			t = t.Elem()
			nullable = true
		case reflect.Map, reflect.Slice, reflect.Array:
			t = t.Elem()
			iterable = true
		case reflect.Func:
			if t.NumOut() == 0 {
				panic("Resolver func must have at least one return value")
			}
			t = t.Out(0)
		}
	}
}

func (my *Builder) getSchema(o *Object) string {
	sb := &strings.Builder{}
	sb.WriteString(my.getDescription(o.Description))

	sb.WriteString("type ")
	sb.WriteString(o.Name)
	sb.WriteString(" {")
	sb.WriteRune('\n')
	for k, v := range o.Resolvers {
		t, n, i := my.getType(reflect.TypeOf(v.Type))
		if t == nil {
			continue
		}
		sb.WriteString(my.getDescription(v.Description))
		sb.WriteString("  ")
		sb.WriteString(k)
		sb.WriteString(": ")
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
	sb.WriteString("}\n")

	return sb.String()
}
