package types

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
