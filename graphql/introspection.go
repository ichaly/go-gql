package graphql

// https://spec.graphql.org/October2021/#sec-Schema-Introspection
type __Schema struct {
	Description      string
	Types            []__Type
	QueryType        *__Type
	MutationType     *__Type
	SubscriptionType *__Type
	Directives       []__Directive
}

type __Type struct {
	Kind        __TypeKind
	Name        *string
	Description *string
	// T_OBJECT and T_INTERFACE only
	Fields func(isDeprecatedArgs) []__Field
	// T_OBJECT only
	Interfaces []__Type
	// T_INTERFACE and T_UNION only
	PossibleTypes func() []__Type
	// T_ENUM only
	EnumValues func(isDeprecatedArgs) []__EnumValue
	// T_INPUT_OBJECT only
	InputFields func() []__InputValue
	// T_NON_NULL and T_LIST only
	OfType *__Type
	// T_SCALAR only
	SpecifiedByURL *string
}

type __Directive struct {
	Name         string
	Description  *string
	Locations    []__DirectiveLocation
	Args         []__InputValue
	IsRepeatable bool
}

type __Field struct {
	Name              string
	Description       *string
	Args              []__InputValue
	Type              __Type
	IsDeprecated      bool
	DeprecationReason *string
}

type __InputValue struct {
	Name         string
	Description  *string
	Type         __Type
	DefaultValue *string
}

type __EnumValue struct {
	Name              string
	Description       string
	IsDeprecated      bool
	DeprecationReason string
}

type isDeprecatedArgs struct {
	IncludeDeprecated bool `json:"includeDeprecated"`
}

type __TypeKind uint8

const (
	TK_SCALAR __TypeKind = iota
	TK_OBJECT
	TK_INTERFACE
	TK_UNION
	TK_ENUM
	TK_INPUT_OBJECT
	TK_LIST
	TK_NON_NULL
)

func (my __TypeKind) String() string {
	switch my {
	case TK_SCALAR:
		return "SCALAR"
	case TK_OBJECT:
		return "OBJECT"
	case TK_INTERFACE:
		return "INTERFACE"
	case TK_UNION:
		return "UNION"
	case TK_ENUM:
		return "ENUM"
	case TK_INPUT_OBJECT:
		return "INPUT_OBJECT"
	case TK_LIST:
		return "LIST"
	case TK_NON_NULL:
		return "NON_NULL"
	default:
		return ""
	}
}

type __DirectiveLocation uint8

const (
	DL_QUERY __DirectiveLocation = iota
	DL_MUTATION
	DL_SUBSCRIPTION
	DL_FIELD
	DL_FRAGMENT_DEFINITION
	DL_FRAGMENT_SPREAD
	DL_INLINE_FRAGMENT
	DL_SCHEMA
	DL_SCALAR
	DL_OBJECT
	DL_FIELD_DEFINITION
	DL_ARGUMENT_DEFINITION
	DL_INTERFACE
	DL_UNION
	DL_ENUM
	DL_ENUM_VALUE
	DL_INPUT_OBJECT
	DL_INPUT_FIELD_DEFINITION
)

func (my __DirectiveLocation) String() string {
	switch my {
	case DL_QUERY:
		return "QUERY"
	case DL_MUTATION:
		return "MUTATION"
	case DL_SUBSCRIPTION:
		return "SUBSCRIPTION"
	case DL_FIELD:
		return "FIELD"
	case DL_FRAGMENT_DEFINITION:
		return "FRAGMENT_DEFINITION"
	case DL_FRAGMENT_SPREAD:
		return "FRAGMENT_SPREAD"
	case DL_INLINE_FRAGMENT:
		return "INLINE_FRAGMENT"
	case DL_SCHEMA:
		return "SCHEMA"
	case DL_SCALAR:
		return "SCALAR"
	case DL_FIELD_DEFINITION:
		return "FIELD_DEFINITION"
	case DL_ARGUMENT_DEFINITION:
		return "ARGUMENT_DEFINITION"
	case DL_INTERFACE:
		return "INTERFACE"
	case DL_UNION:
		return "UNION"
	case DL_OBJECT:
		return "OBJECT"
	case DL_ENUM:
		return "ENUM"
	case DL_ENUM_VALUE:
		return "ENUM_VALUE"
	case DL_INPUT_OBJECT:
		return "INPUT_OBJECT"
	case DL_INPUT_FIELD_DEFINITION:
		return "INPUT_FIELD_DEFINITION"
	default:
		return ""
	}
}
