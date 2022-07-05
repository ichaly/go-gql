package graphql

import (
	"fmt"
	"reflect"
)

type Object struct {
	key         string // Optional, defaults 'ID'.
	Name        string // Optional, defaults to Value's name.
	Description string
	Value       interface{}
	Fields      map[string]*Object
}

type FieldOption func(*Object)

func (my *Object) Field(name string, value interface{}, options ...FieldOption) {
	if my.Fields == nil {
		my.Fields = make(map[string]*Object)
	}
	if _, ok := my.Fields[name]; ok {
		panic("duplicate Field")
	}
	obj := &Object{Value: value}
	for _, o := range options {
		o(obj)
	}
	my.Fields[name] = obj
}

func (my *Object) String() (result string) {
	t := reflect.TypeOf(my.Value)
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
		if my.Name == "ID" || my.Name == my.key {
			result = fmt.Sprintf(result, "ID")
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
			result = fmt.Sprintf(result, my.Name)
			return
		}
		return
	}
}
