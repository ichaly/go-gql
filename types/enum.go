package types

type Enum struct {
	Type       string
	Values     []string
	ReverseMap map[interface{}]string
	Map        map[string]interface{}
}
