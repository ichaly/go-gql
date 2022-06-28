package types

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
