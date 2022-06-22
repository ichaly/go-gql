package graphql

type Object struct {
	Name        string // Optional, defaults to Type's name.
	Description string
	Type        interface{}
	key         string
	ServiceName string
	Methods     map[string]*Method
}

func (s *Object) Key(f string) {
	s.key = f
}

type FieldOption interface {
	apply(*Method)
}

func (s *Object) Field(name string, f interface{}, options ...FieldOption) {
	if s.Methods == nil {
		s.Methods = make(map[string]*Method)
	}

	m := &Method{Fn: f}
	for _, o := range options {
		o.apply(m)
	}

	if _, ok := s.Methods[name]; ok {
		panic("duplicate Method")
	}
	s.Methods[name] = m
}

type Method struct {
	Fn interface{}
}

type Enum struct {
	Type       string
	Values     []string
	ReverseMap map[interface{}]string
	Map        map[string]interface{}
}
