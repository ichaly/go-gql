package types

import (
	"encoding/json"
	"fmt"
	"github.com/ichaly/go-gql/executor"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"net/http"
	"time"
)

type (
	Transport interface {
		Supports(r *http.Request) bool
		Do(w http.ResponseWriter, r *http.Request, exec *executor.Executor)
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

func SendErrorf(w http.ResponseWriter, code int, format string, args ...interface{}) {
	SendError(w, code, &gqlerror.Error{Message: fmt.Sprintf(format, args...)})
}

func SendError(w http.ResponseWriter, code int, errors ...*gqlerror.Error) {
	w.WriteHeader(code)
	b, err := json.Marshal(&GqlResult{Errors: errors})
	if err != nil {
		panic(err)
	}
	_, _ = w.Write(b)
}
