package graphql

import (
	"fmt"
	"github.com/ichaly/go-gql/internal/introspection"
	"reflect"
)

type SchemaBuilder struct {
	Name    string
	enums   map[reflect.Type]*Enum
	objects map[reflect.Type]*Object
	types   map[reflect.Type]*introspection.Type
}

type BuildOption interface {
	apply(*SchemaBuilder)
}

func NewSchemaBuilder(options ...BuildOption) *SchemaBuilder {
	builder := &SchemaBuilder{
		objects: make(map[reflect.Type]*Object),
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
	if object, ok := my.objects[reflect.TypeOf(typ)]; ok {
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
	my.objects[reflect.TypeOf(typ)] = object

	for _, o := range options {
		o.apply(my, object)
	}

	return object
}

func (my *SchemaBuilder) Enum(val interface{}, enumMap interface{}) {
	typ := reflect.TypeOf(val)
	if my.enums == nil {
		my.enums = make(map[reflect.Type]*Enum)
	}
	eMap, rMap := getEnumMap(enumMap, typ)
	my.enums[typ] = &Enum{Map: eMap, ReverseMap: rMap}
}

func getEnumMap(enumMap interface{}, typ reflect.Type) (map[string]interface{}, map[interface{}]string) {
	rMap := make(map[interface{}]string)
	eMap := make(map[string]interface{})
	v := reflect.ValueOf(enumMap)
	if v.Kind() == reflect.Map {
		for _, key := range v.MapKeys() {
			val := v.MapIndex(key)
			valInterface := val.Interface()
			if reflect.TypeOf(valInterface).Kind() != typ.Kind() {
				panic("types are not equal")
			}
			if key.Kind() == reflect.String {
				mapVal := reflect.ValueOf(valInterface).Convert(typ)
				eMap[key.String()] = mapVal.Interface()
			} else {
				panic("keys are not strings")
			}
		}
	} else {
		panic("enum function not passed a map")
	}

	for key, val := range eMap {
		rMap[val] = key
	}
	return eMap, rMap
}

func (my *SchemaBuilder) Build() (*introspection.Schema, error) {
	my.Object("Query", query{})
	my.Object("Mutation", mutation{})
	return &introspection.Schema{
		QueryType:    nil,
		MutationType: nil,
	}, nil
}

func (my *SchemaBuilder) MustBuild() *introspection.Schema {
	schema, err := my.Build()
	if err != nil {
		panic(err)
	}
	return schema
}

func (my *SchemaBuilder) getType(nodeType reflect.Type) (*introspection.Type, error) {
	// Structs
	if nodeType.Kind() == reflect.Struct {
		if err := my.buildStruct(nodeType); err != nil {
			return nil, err
		}
		return &introspection.Type{Kind: introspection.NON_NULL, OfType: my.types[nodeType]}, nil
	}
	if nodeType.Kind() == reflect.Ptr && nodeType.Elem().Kind() == reflect.Struct {
		if err := my.buildStruct(nodeType.Elem()); err != nil {
			return nil, err
		}
		return my.types[nodeType.Elem()], nil
	}
	if nodeType.Kind() == reflect.Slice {
		elementType, err := my.getType(nodeType.Elem())
		if err != nil {
			return nil, err
		}
		return &introspection.Type{Kind: introspection.NON_NULL, OfType: &introspection.Type{Kind: introspection.LIST, OfType: elementType}}, nil
	}
	return nil, fmt.Errorf("bad type %s: should be a scalar, slice, or struct type", nodeType)
}

func (sb *SchemaBuilder) buildStruct(typ reflect.Type) error {
	return nil
}
