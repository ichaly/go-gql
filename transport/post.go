package transport

import (
	"github.com/ichaly/go-gql/types"
	"github.com/ichaly/go-gql/util"
	"mime"
	"net/http"
)

// POST implements the POST side of the default HTTP transport
// defined in https://github.com/APIs-guru/graphql-over-http#post
type POST struct{}

func (h POST) Supports(r *http.Request) bool {
	if r.Header.Get("Upgrade") != "" {
		return false
	}

	mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil {
		return false
	}

	return r.Method == "POST" && mediaType == "application/json"
}

func (h POST) Do(w http.ResponseWriter, r *http.Request, exec *types.Executor) {
	w.Header().Set("Content-Type", "application/json")

	start := util.Now()
	var params *types.GqlRequest
	if err := util.ReadJson(r.Body, &params); err != nil {
		types.SendErrorf(w, http.StatusBadRequest, "json body could not be decoded: "+err.Error())
		return
	}
	params.Headers = r.Header
	params.ReadTime = types.TraceTiming{
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
	//responses, ctx := exec.DispatchOperation(r.Context(), rc)
	//writeJson(w, responses(ctx))
}
