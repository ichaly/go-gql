package graphql

import "fmt"

type IType interface {
	Alias() string
}

type NonNull struct {
	Type IType
}

func (my *NonNull) Alias() string {
	return fmt.Sprintf("%s!", my.Type)
}

type Enum struct {
	Type       string
	Values     []string
	ReverseMap map[interface{}]string
}

func (my *Enum) Alias() string {
	return my.Type
}

type Scalar struct {
	Type      string
	UnWrapper func(interface{}) (interface{}, error)
}

func (s *Scalar) Alias() string {
	return s.Type
}
