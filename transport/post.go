package transport

import (
	"github.com/ichaly/go-gql/executor"
	"github.com/ichaly/go-gql/util"
	"log"
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

func (h POST) Do(w http.ResponseWriter, r *http.Request, exec *executor.Executor) {
	w.Header().Set("Content-Type", "application/json")

	start := util.Now()
	var params *executor.GqlRequest
	if err := util.ReadJson(r.Body, &params); err != nil {
		util.SendErrorf(w, http.StatusBadRequest, "json body could not be decoded: "+err.Error())
		return
	}
	params.Headers = r.Header
	params.ReadTime = executor.TraceTiming{
		Start: start,
		End:   util.Now(),
	}
	res := exec.Exec(r.Context(), params)
	log.Fatalf("%v", res)
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
