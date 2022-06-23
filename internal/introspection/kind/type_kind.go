package kind

type TypeKind uint8

const (
	SCALAR TypeKind = iota
	OBJECT
	INTERFACE
	UNION
	ENUM
	INPUT_OBJECT
	LIST
	NON_NULL
)

func (my TypeKind) String() string {
	switch my {
	case SCALAR:
		return "SCALAR"
	case OBJECT:
		return "OBJECT"
	case INTERFACE:
		return "INTERFACE"
	case UNION:
		return "UNION"
	case ENUM:
		return "ENUM"
	case INPUT_OBJECT:
		return "INPUT_OBJECT"
	case LIST:
		return "LIST"
	case NON_NULL:
		return "NON_NULL"
	default:
		return ""
	}
}

func (TypeKind) Alias() string {
	return "__TypeKind"
}
