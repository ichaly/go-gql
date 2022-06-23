package graphql

import (
	h "github.com/ichaly/go-gql/graphql/internal"
	"github.com/ichaly/go-gql/internal/introspection"
)

var (
	scalarAny = introspection.Type{
		Kind:        introspection.SCALAR,
		Name:        h.StrPtr("Any"),
		Description: h.StrPtr("The `Any` scalar type represents interface{}."),
	}
	scalarBoolean = introspection.Type{
		Kind:        introspection.SCALAR,
		Name:        h.StrPtr("Boolean"),
		Description: h.StrPtr("The `Boolean` scalar type represents `true` or `false`."),
	}
	scalarInt = introspection.Type{
		Kind:        introspection.SCALAR,
		Name:        h.StrPtr("Int"),
		Description: h.StrPtr("The Int scalar type represents a signed 32‐bit numeric non‐fractional value."),
	}
	scalarFloat = introspection.Type{
		Kind:        introspection.SCALAR,
		Name:        h.StrPtr("Float"),
		Description: h.StrPtr("The Float scalar type represents signed double‐precision fractional values as specified by IEEE 754."),
	}
	scalarString = introspection.Type{
		Kind:        introspection.SCALAR,
		Name:        h.StrPtr("String"),
		Description: h.StrPtr("The `String` scalar type represents textual data, represented as UTF-8 character sequences. The String type is most often used by GraphQL to represent free-form human-readable text."),
	}
	scalarID = introspection.Type{
		Kind:        introspection.SCALAR,
		Name:        h.StrPtr("ID"),
		Description: h.StrPtr("The ID scalar type represents a unique identifier, often used to refetch an object or as the key for a cache"),
	}
	scalarFile = introspection.Type{
		Kind:           introspection.SCALAR,
		Name:           h.StrPtr("File"),
		Description:    h.StrPtr("The File scalar type references to a multipart file, often used to upload files to the server. Expects a string with the form file field name"),
		SpecifiedByURL: h.StrPtr("https://github.com/mjarkk/yarql#file-upload"),
	}
	scalarTime = introspection.Type{
		Kind:           introspection.SCALAR,
		Name:           h.StrPtr("Time"),
		Description:    h.StrPtr("The Time scalar type references to a ISO 8601 date+time, often used to insert and/or view dates. Expects a string with the ISO 8601 format"),
		SpecifiedByURL: h.StrPtr("https://en.wikipedia.org/wiki/ISO_8601"),
	}
)

var scalars = map[string]introspection.Type{
	"Any":     scalarAny,
	"Boolean": scalarBoolean,
	"Int":     scalarInt,
	"Float":   scalarFloat,
	"String":  scalarString,
	"ID":      scalarID,
	"File":    scalarFile,
	"Time":    scalarTime,
}
