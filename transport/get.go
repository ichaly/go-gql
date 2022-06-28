package transport

import (
	"github.com/ichaly/go-gql"
	"github.com/ichaly/go-gql/types"
	"github.com/ichaly/go-gql/util"
	"net/http"
	"strings"
)

// GET implements the GET side of the default HTTP transport
// defined in https://github.com/APIs-guru/graphql-over-http#get
type GET struct{}

var _ graphql.Transport = GET{}

func (h GET) Supports(r *http.Request) bool {
	if r.Header.Get("Upgrade") != "" {
		return false
	}

	return r.Method == "GET"
}

func (h GET) Do(w http.ResponseWriter, r *http.Request, exec *graphql.Executor) {
	w.Header().Set("Content-Type", "application/json")

	start := util.Now()
	params := &graphql.GqlRequest{
		Query:         r.URL.Query().Get("query"),
		OperationName: r.URL.Query().Get("operationName"),
		Headers:       r.Header,
	}
	if variables := r.URL.Query().Get("variables"); variables != "" {
		if err := util.ReadJson(strings.NewReader(variables), &params.Variables); err != nil {
			types.SendErrorf(w, http.StatusBadRequest, "variables could not be decoded")
			return
		}
	}
	if extensions := r.URL.Query().Get("extensions"); extensions != "" {
		if err := util.ReadJson(strings.NewReader(extensions), &params.Extensions); err != nil {
			types.SendErrorf(w, http.StatusBadRequest, "extensions could not be decoded")
			return
		}
	}
	params.ReadTime = graphql.TraceTiming{
		Start: start,
		End:   util.Now(),
	}

	//rc, err := exec.CreateOperationContext(r.Context(), params)
	//if err != nil {
	//	w.WriteHeader(statusFor(err))
	//	resp := exec.DispatchError(graphql.WithOperationContext(r.Context(), rc), err)
	//	writeJson(w, resp)
	//	return
	//}
	//op := rc.Doc.Operations.ForName(rc.OperationName)
	//if op.Operation != ast.Query {
	//	w.WriteHeader(http.StatusNotAcceptable)
	//	writeJsonError(w, "GET requests only allow query operations")
	//	return
	//}
	//
	//responses, ctx := exec.DispatchOperation(r.Context(), rc)
	//writeJson(w, responses(ctx))
}