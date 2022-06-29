package executor

import (
	"context"
	"encoding/json"
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"github.com/vektah/gqlparser/v2/parser"
	"github.com/vektah/gqlparser/v2/validator"
	"net/http"
	"time"
)

type (
	Executor struct {
		schema   *ast.Schema
		query    interface{}
		mutation interface{}
	}
	TraceTiming struct {
		Start time.Time
		End   time.Time
	}
	GqlRequest struct {
		Query         string                 `json:"query"`
		OperationName string                 `json:"operationName"`
		Variables     map[string]interface{} `json:"variables"`
		Extensions    map[string]interface{} `json:"extensions"`
		Headers       http.Header            `json:"headers"`

		ReadTime TraceTiming `json:"-"`
	}
	GqlResult struct {
		Data       json.RawMessage        `json:"data,omitempty"`
		Errors     gqlerror.List          `json:"errors,omitempty"`
		Extensions map[string]interface{} `json:"extensions,omitempty"`
	}
)

func (my *Executor) Exec(
	ctx context.Context, req *GqlRequest,
) (r GqlResult) {
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
