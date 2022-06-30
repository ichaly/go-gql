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
	return my.Object(&query{}, WithName("Query"))
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

func (my *Builder) Object(obj interface{}, options ...objectOption) *Object {
	typ := reflect.TypeOf(obj)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		obj = reflect.ValueOf(obj).Elem().Interface()
	}
	if object, ok := my.objects[typ]; ok {
		if reflect.TypeOf(object.Type) != typ {
			panic("re-registered object with different type")
		}
		return object
	}
	object := &Object{
		Name: typ.Name(),
		Type: obj,
	}
	my.objects[typ] = object

	for _, o := range options {
		o(my, object)
	}

	for i := 0; i < typ.NumField(); i++ {
		sf := typ.Field(i)
		object.Field(sf.Name, reflect.New(sf.Type).Elem().Interface())
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
	sb.WriteString("scalar Any\n\n")

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

var (
	scalars = map[reflect.Kind]*Object{
		reflect.Int: {
			Name:        "Int",
			Description: "The Int scalar type represents a signed 32‐bit numeric non‐fractional value.",
		},
		reflect.Float64: {
			Name:        "Float",
			Description: "The Float scalar type represents signed double‐precision fractional values as specified by IEEE 754.",
		},
		reflect.Bool: {
			Name:        "Boolean",
			Description: "The `Boolean` scalar type represents `true` or `false`.",
		},
		reflect.String: {
			Name:        "String",
			Description: "The `String` scalar type represents textual data, represented as UTF-8 character sequences. The String type is most often used by GraphQL to represent free-form human-readable text.",
		},
		reflect.Interface: {
			Name:        "Any",
			Description: "The `Any` scalar type represents interface{}.",
		},
	}
)

func (my *Builder) getType(t reflect.Type) (result *Object, nullable bool, iterable bool) {
	for {
		switch t.Kind() {
		case reflect.Bool, reflect.String, reflect.Interface:
			result = scalars[t.Kind()]
			return
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			result = scalars[reflect.Int]
			return
		case reflect.Float32, reflect.Float64:
			result = scalars[reflect.Float64]
			return
		case reflect.Struct:
			result = my.Object(reflect.New(t).Elem().Interface())
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
		o, n, i := my.getType(reflect.TypeOf(v.Type))
		if o == nil {
			continue
		}
		sb.WriteString(my.getDescription(v.Description))
		sb.WriteString("  ")
		sb.WriteString(k)
		sb.WriteString(": ")
		if i {
			sb.WriteString("[")
		}
		if k == "ID" || o.key == k {
			sb.WriteString("ID")
		} else {
			sb.WriteString(o.Name)
		}
		if n {
			sb.WriteString("!")
		}
		if i {
			sb.WriteString("]!")
		}
		sb.WriteRune('\n')
	}
	sb.WriteString("}\n\n")

	return sb.String()
}
