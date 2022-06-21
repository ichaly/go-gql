package graphql

type NonNull struct {
	Type Type
}

type Enum struct {
	Type       string
	Values     []string
	ReverseMap map[interface{}]string
}

func (e *Enum) isType() {}

func (e *Enum) String() string {
	return e.Type
}

func (e *Enum) enumValues() []string {
	return e.Values
}
