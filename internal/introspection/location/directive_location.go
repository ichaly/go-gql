package location

type DirectiveLocation uint8

const (
	QUERY DirectiveLocation = iota
	MUTATION
	SUBSCRIPTION
	FIELD
	FRAGMENT_DEFINITION
	FRAGMENT_SPREAD
	INLINE_FRAGMENT
	VARIABLE_DEFINITION
	SCHEMA
	SCALAR
	OBJECT
	FIELD_DEFINITION
	ARGUMENT_DEFINITION
	INTERFACE
	UNION
	ENUM
	ENUM_VALUE
	INPUT_OBJECT
	INPUT_FIELD_DEFINITION
)

func (my DirectiveLocation) String() string {
	switch my {
	case QUERY:
		return "QUERY"
	case MUTATION:
		return "MUTATION"
	case SUBSCRIPTION:
		return "SUBSCRIPTION"
	case FIELD:
		return "FIELD"
	case FRAGMENT_DEFINITION:
		return "FRAGMENT_DEFINITION"
	case FRAGMENT_SPREAD:
		return "FRAGMENT_SPREAD"
	case INLINE_FRAGMENT:
		return "INLINE_FRAGMENT"
	case VARIABLE_DEFINITION:
		return "VARIABLE_DEFINITION"
	case SCHEMA:
		return "SCHEMA"
	case SCALAR:
		return "SCALAR"
	case FIELD_DEFINITION:
		return "FIELD_DEFINITION"
	case ARGUMENT_DEFINITION:
		return "ARGUMENT_DEFINITION"
	case INTERFACE:
		return "INTERFACE"
	case UNION:
		return "UNION"
	case OBJECT:
		return "OBJECT"
	case ENUM:
		return "ENUM"
	case ENUM_VALUE:
		return "ENUM_VALUE"
	case INPUT_OBJECT:
		return "INPUT_OBJECT"
	case INPUT_FIELD_DEFINITION:
		return "INPUT_FIELD_DEFINITION"
	default:
		return ""
	}
}

func (DirectiveLocation) Alias() string {
	return "__DirectiveLocation"
}
