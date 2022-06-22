package graphql

import (
	"fmt"
	"github.com/ichaly/go-gql/graphql/internal"
	"reflect"
	"time"
)

type SchemaBuilder struct {
	Name    string
	objects map[reflect.Type]*Object
	enums   map[reflect.Type]*EnumMapping
}

type EnumMapping struct {
	Map        map[string]interface{}
	ReverseMap map[interface{}]string
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
		my.enums = make(map[reflect.Type]*EnumMapping)
	}
	eMap, rMap := getEnumMap(enumMap, typ)
	my.enums[typ] = &EnumMapping{Map: eMap, ReverseMap: rMap}
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

func (my *SchemaBuilder) Build() (*Schema, error) {
	my.Object("Query", query{})
	my.Object("Mutation", mutation{})
	for _, object := range my.objects {
		typ := reflect.TypeOf(object.Type)
		if typ.Kind() != reflect.Struct {
			return nil, fmt.Errorf("object.IType should be a struct, not %my", typ.String())
		}

		if _, ok := my.objects[typ]; ok {
			return nil, fmt.Errorf("duplicate object for %my", typ.String())
		}

		my.objects[typ] = object
	}
	queryTyp, err := my.getType(reflect.TypeOf(&query{}))
	if err != nil {
		return nil, err
	}
	mutationTyp, err := my.getType(reflect.TypeOf(&mutation{}))
	if err != nil {
		return nil, err
	}
	return &Schema{
		QueryType:    &queryTyp,
		MutationType: &mutationTyp,
	}, nil
}

func (my *SchemaBuilder) MustBuild() *Schema {
	schema, err := my.Build()
	if err != nil {
		panic(err)
	}
	return schema
}

func (my *SchemaBuilder) getEnum(typ reflect.Type) (string, []string, bool) {
	if my.enums[typ] != nil {
		var values []string
		for mapping := range my.enums[typ].Map {
			values = append(values, mapping)
		}
		return typ.Name(), values, true
	}
	return "", nil, false
}

var scalars = map[reflect.Type]string{
	reflect.TypeOf(bool(false)): "bool",
	reflect.TypeOf(int(0)):      "int",
	reflect.TypeOf(int8(0)):     "int8",
	reflect.TypeOf(int16(0)):    "int16",
	reflect.TypeOf(int32(0)):    "int32",
	reflect.TypeOf(int64(0)):    "int64",
	reflect.TypeOf(uint(0)):     "uint",
	reflect.TypeOf(uint8(0)):    "uint8",
	reflect.TypeOf(uint16(0)):   "uint16",
	reflect.TypeOf(uint32(0)):   "uint32",
	reflect.TypeOf(uint64(0)):   "uint64",
	reflect.TypeOf(float32(0)):  "float32",
	reflect.TypeOf(float64(0)):  "float64",
	reflect.TypeOf(string("")):  "string",
	reflect.TypeOf(time.Time{}): "Time",
	reflect.TypeOf([]byte{}):    "bytes",
}

func getScalar(typ reflect.Type) (string, bool) {
	for match, name := range scalars {
		if internal.TypesIdenticalOrScalarAliases(match, typ) {
			return name, true
		}
	}
	return "", false
}

func (my *SchemaBuilder) getType(nodeType reflect.Type) (IType, error) {
	if typeName, values, ok := my.getEnum(nodeType); ok {
		return &NonNull{
			Type: &Enum{
				Type:       typeName,
				Values:     values,
				ReverseMap: my.enums[nodeType].ReverseMap,
			},
		}, nil
	}

	if typeName, ok := getScalar(nodeType); ok {
		return &NonNull{Type: &Scalar{Type: typeName}}, nil
	}
	if nodeType.Kind() == reflect.Ptr {
		if typeName, ok := getScalar(nodeType.Elem()); ok {
			return &Scalar{Type: typeName}, nil // XXX: prefix typ with "*"
		}
	}

	if nodeType.Implements(textMarshalerType) {
		return my.getTextMarshalerType(nodeType)
	}

	// Structs
	if nodeType.Kind() == reflect.Struct {
		if err := my.buildStruct(nodeType); err != nil {
			return nil, err
		}
		return &NonNull{Type: my.types[nodeType]}, nil
	}
	if nodeType.Kind() == reflect.Ptr && nodeType.Elem().Kind() == reflect.Struct {
		if err := my.buildStruct(nodeType.Elem()); err != nil {
			return nil, err
		}
		return my.types[nodeType.Elem()], nil
	}

	switch nodeType.Kind() {
	case reflect.Slice:
		elementType, err := my.getType(nodeType.Elem())
		if err != nil {
			return nil, err
		}

		// Wrap all slice elements in NonNull.
		if _, ok := elementType.(*NonNull); !ok {
			elementType = &NonNull{Type: elementType}
		}

		return &NonNull{Type: &List{Type: elementType}}, nil

	default:
		return nil, fmt.Errorf("bad type %s: should be a scalar, slice, or struct type", nodeType)
	}
}
