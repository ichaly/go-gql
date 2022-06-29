package introspection

import (
	"github.com/ichaly/go-gql/internal/introspection/kind"
	"github.com/ichaly/go-gql/internal/introspection/location"
)

// Node https://spec.graphql.org/October2021/#sec-Schema-Introspection
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
	Kind        kind.TypeKind
	Name        *string
	Description *string
	// OBJECT and INTERFACE only
	Fields func(isDeprecatedArgs) []Field
	// OBJECT only
	Interfaces []Type
	// INTERFACE and UNION only
	PossibleTypes func() []Type
	// ENUM only
	EnumValues func(isDeprecatedArgs) []EnumValue
	// INPUOBJECT only
	InputFields func() []InputValue
	// NON_NULL and LIST only
	OfType *Type
	// SCALAR only
	SpecifiedByURL *string
}

func (Type) Alias() string {
	return "__Type"
}

type Directive struct {
	Name         string
	Description  *string
	Locations    []location.DirectiveLocation
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
