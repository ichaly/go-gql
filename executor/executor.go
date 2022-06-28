package executor

import (
	"context"
	"github.com/ichaly/go-gql/types"
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"github.com/vektah/gqlparser/v2/parser"
	"github.com/vektah/gqlparser/v2/validator"
)

type Executor struct {
	schema   *ast.Schema
	query    interface{}
	mutation interface{}
}

func (my *Executor) CreateOperationContext(
	ctx context.Context, req *types.GqlRequest,
) (r types.GqlResult) {
	doc, err := parser.ParseQuery(&ast.Source{Input: req.Query})
	if err != nil {
		r.Errors = append(r.Errors, err)
		return
	}
	if len(doc.Operations) == 0 {
		r.Errors = append(r.Errors, gqlerror.Errorf("no operation provided"))
		return
	}
	r.Errors = validator.Validate(my.schema, doc)
	if r.Errors != nil {
		return
	}
	op := doc.Operations.ForName(req.OperationName)
	if op == nil {
		r.Errors = append(r.Errors, gqlerror.Errorf("operation %s not found", req.OperationName))
		return
	}
	_, err = validator.VariableValues(my.schema, op, req.Variables)
	if err != nil {
		r.Errors = append(r.Errors, err)
		return
	}

	switch op.Operation {
	case ast.Query:
	case ast.Mutation:
	case ast.Subscription:
	default:
	}
	return
}
