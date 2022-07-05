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
	Key         string // Optional, defaults 'ID'.
	Name        string // Optional, defaults to Value's name.
	Description string
	Value       interface{}
	Fields      map[string]*Object
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

type ObjectOption func(*Builder, *Object)

func WithName(name string) ObjectOption {
	return func(b *Builder, o *Object) {
		o.Name = name
	}
}

func (my *Builder) Object(value interface{}, options ...ObjectOption) *Object {
	typ := reflect.TypeOf(value)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		value = reflect.ValueOf(value).Elem().Interface()
	}
	if typ.Kind() != reflect.Struct {
		panic("object only can be set a struct or struct pointer")
	}
	if object, ok := my.objects[typ]; ok {
		if reflect.TypeOf(object.Value) != typ {
			panic("re-registered object with different type")
		}
		return object
	}
	object := &Object{
		Name:  typ.Name(),
		Value: value,
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

type FieldOption interface {
	apply(*Object)
}

func (s *Object) Field(name string, value interface{}, options ...FieldOption) {
	if s.Fields == nil {
		s.Fields = make(map[string]*Object)
	}
	if _, ok := s.Fields[name]; ok {
		panic("duplicate Field")
	}
	obj := &Object{Value: value}
	for _, o := range options {
		o.apply(obj)
	}
	s.Fields[name] = obj
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
	sb.WriteString("scalar Time\n\n")

	for _, o := range my.objects {
		sb.WriteString(my.getSchema(o))
	}

	fmt.Println(sb.String())
	s, e := gqlparser.LoadSchema(&ast.Source{
		Name:  "schema",
		Input: sb.String(),
	})
	return s, e
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

func (my *Builder) getName(obj *Object) (result string) {
	t := reflect.TypeOf(obj.Value)
	result = "%s"
	for {
		k := t.Kind()
		switch k {
		case reflect.Ptr:
			t = t.Elem()
			result = fmt.Sprintf(result, "%s!")
			continue
		case reflect.Map, reflect.Slice, reflect.Array:
			t = t.Elem()
			result = fmt.Sprintf(result, "[%s]!")
			continue
		case reflect.Func:
			if t.NumOut() == 0 {
				panic("Resolver func must have at least one return value")
			}
			t = t.Out(0)
			continue
		}
		if len(obj.Name) > 0 && (obj.Name == "ID" || obj.Name == obj.Key) {
			result = fmt.Sprintf(result, "ID")
			return
		}
		if t.Name() == "Time" {
			result = fmt.Sprintf(result, "Time")
			return
		}
		switch k {
		case reflect.Bool:
			result = fmt.Sprintf(result, "Boolean")
			return
		case reflect.String:
			result = fmt.Sprintf(result, "String")
			return
		case reflect.Interface:
			result = fmt.Sprintf(result, "Any")
			return
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			result = fmt.Sprintf(result, "Int")
			return
		case reflect.Float32, reflect.Float64:
			result = fmt.Sprintf(result, "Float")
			return
		case reflect.Struct:
			result = fmt.Sprintf(result, t.Name())
			my.Object(reflect.New(t).Elem().Interface())
			return
		}
		result = ""
		return
	}
}

func (my *Builder) getSchema(o *Object) string {
	sb := &strings.Builder{}
	sb.WriteString(my.getDescription(o.Description))

	sb.WriteString("type ")
	sb.WriteString(o.Name)
	sb.WriteString(" {")
	sb.WriteRune('\n')
	for k, v := range o.Fields {
		name := my.getName(v)
		if len(name) == 0 {
			continue
		}
		sb.WriteString(my.getDescription(v.Description))
		sb.WriteString("  ")
		sb.WriteString(k)
		sb.WriteString(": ")
		sb.WriteString(name)
		sb.WriteRune('\n')
	}
	sb.WriteString("}\n\n")

	return sb.String()
}
