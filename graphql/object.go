package graphql

type Object struct {
	Name        string // Optional, defaults to Type's name.
	Description string
	Type        interface{}
	key         string
	ServiceName string
	Methods     map[string]*method
}

func (s *Object) Key(f string) {
	s.key = f
}

type FieldOption interface {
	apply(*method)
}

func (s *Object) Field(name string, f interface{}, options ...FieldOption) {
	if s.Methods == nil {
		s.Methods = make(map[string]*method)
	}

	m := &method{Fn: f}
	for _, o := range options {
		o.apply(m)
	}

	if _, ok := s.Methods[name]; ok {
		panic("duplicate method")
	}
	s.Methods[name] = m
}
