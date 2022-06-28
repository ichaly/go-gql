package introspection

import (
	"github.com/ichaly/go-gql/internal/introspection/kind"
	"github.com/ichaly/go-gql/util"
)

var (
	scalarAny = Type{
		Kind:        kind.SCALAR,
		Name:        util.StrPtr("Any"),
		Description: util.StrPtr("The `Any` scalar type represents interface{}."),
	}
	scalarBoolean = Type{
		Kind:        kind.SCALAR,
		Name:        util.StrPtr("Boolean"),
		Description: util.StrPtr("The `Boolean` scalar type represents `true` or `false`."),
	}
	scalarInt = Type{
		Kind:        kind.SCALAR,
		Name:        util.StrPtr("Int"),
		Description: util.StrPtr("The Int scalar type represents a signed 32‐bit numeric non‐fractional value."),
	}
	scalarFloat = Type{
		Kind:        kind.SCALAR,
		Name:        util.StrPtr("Float"),
		Description: util.StrPtr("The Float scalar type represents signed double‐precision fractional values as specified by IEEE 754."),
	}
	scalarString = Type{
		Kind:        kind.SCALAR,
		Name:        util.StrPtr("String"),
		Description: util.StrPtr("The `String` scalar type represents textual data, represented as UTF-8 character sequences. The String type is most often used by GraphQL to represent free-form human-readable text."),
	}
	scalarID = Type{
		Kind:        kind.SCALAR,
		Name:        util.StrPtr("ID"),
		Description: util.StrPtr("The ID scalar type represents a unique identifier, often used to refetch an object or as the key for a cache"),
	}
	scalarFile = Type{
		Kind:           kind.SCALAR,
		Name:           util.StrPtr("File"),
		Description:    util.StrPtr("The File scalar type references to a multipart file, often used to upload files to the server. Expects a string with the form file field name"),
		SpecifiedByURL: util.StrPtr("https://github.com/mjarkk/yarql#file-upload"),
	}
	scalarTime = Type{
		Kind:           kind.SCALAR,
		Name:           util.StrPtr("Time"),
		Description:    util.StrPtr("The Time scalar type references to a ISO 8601 date+time, often used to insert and/or view dates. Expects a string with the ISO 8601 format"),
		SpecifiedByURL: util.StrPtr("https://en.wikipedia.org/wiki/ISO_8601"),
	}
)

var scalars = map[string]Type{
	"Any":     scalarAny,
	"Boolean": scalarBoolean,
	"Int":     scalarInt,
	"Float":   scalarFloat,
	"String":  scalarString,
	"ID":      scalarID,
	"File":    scalarFile,
	"Time":    scalarTime,
}
