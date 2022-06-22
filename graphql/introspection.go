package graphql

// https://spec.graphql.org/October2021/#sec-Schema-Introspection
type Schema struct {
	Description      string
	Types            []Type
	QueryType        *Type
	MutationType     *Type
	SubscriptionType *Type
	Directives       []Directive
}

func (Schema) Alias() string {
	return "__Schema"
}

type Type struct {
	Kind        TypeKind
	Name        *string
	Description *string
	// T_OBJECT and T_INTERFACE only
	Fields func(isDeprecatedArgs) []Field
	// T_OBJECT only
	Interfaces []Type
	// T_INTERFACE and T_UNION only
	PossibleTypes func() []Type
	// T_ENUM only
	EnumValues func(isDeprecatedArgs) []EnumValue
	// T_INPUT_OBJECT only
	InputFields func() []InputValue
	// T_NON_NULL and T_LIST only
	OfType *Type
	// T_SCALAR only
	SpecifiedByURL *string
}

func (Type) Alias() string {
	return "__Type"
}

type Directive struct {
	Name         string
	Description  *string
	Locations    []DirectiveLocation
	Args         []InputValue
	IsRepeatable bool
}

func (Directive) Alias() string {
	return "__Directive"
}

type Field struct {
	Name              string
	Description       *string
	Args              []InputValue
	Type              Type
	IsDeprecated      bool
	DeprecationReason *string
}

func (Field) Alias() string {
	return "__Field"
}

type InputValue struct {
	Name         string
	Description  *string
	Type         Type
	DefaultValue *string
}

func (InputValue) Alias() string {
	return "__InputValue"
}

type EnumValue struct {
	Name              string
	Description       string
	IsDeprecated      bool
	DeprecationReason string
}

func (EnumValue) Alias() string {
	return "__EnumValue"
}

type isDeprecatedArgs struct {
	IncludeDeprecated bool `json:"includeDeprecated"`
}

type TypeKind uint8

const (
	TK_SCALAR TypeKind = iota
	TK_OBJECT
	TK_INTERFACE
	TK_UNION
	TK_ENUM
	TK_INPUT_OBJECT
	TK_LIST
	TK_NON_NULL
)

func (my TypeKind) String() string {
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

func (TypeKind) Alias() string {
	return "__TypeKind"
}

type DirectiveLocation uint8

const (
	DL_QUERY DirectiveLocation = iota
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

func (my DirectiveLocation) String() string {
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

func (DirectiveLocation) Alias() string {
	return "__DirectiveLocation"
}
